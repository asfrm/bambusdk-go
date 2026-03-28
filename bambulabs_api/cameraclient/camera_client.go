package cameraclient

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/bambulabs_api/go/bambulabs_api/logger"
)

// PrinterCamera handles camera stream from the printer
type PrinterCamera struct {
	username   string
	accessCode string
	hostname   string
	port       int

	mu        sync.RWMutex
	thread    chan struct{}
	lastFrame []byte
	alive     bool
	stopChan  chan struct{}
}

// NewPrinterCamera creates a new camera client
func NewPrinterCamera(hostname, accessCode string, port int, username string) *PrinterCamera {
	if port == 0 {
		port = 6000
	}
	if username == "" {
		username = "bblp"
	}

	return &PrinterCamera{
		username:   username,
		accessCode: accessCode,
		hostname:   hostname,
		port:       port,
		alive:      false,
		stopChan:   make(chan struct{}),
	}
}

// Start starts the camera client
func (c *PrinterCamera) Start() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.alive {
		return false
	}

	c.alive = true
	c.thread = make(chan struct{})
	go c.retriever()

	logger.Info.Println("Starting camera thread")
	return true
}

// Stop stops the camera client
func (c *PrinterCamera) Stop() {
	c.mu.Lock()
	if !c.alive {
		c.mu.Unlock()
		return
	}
	c.alive = false
	c.mu.Unlock()

	close(c.stopChan)

	if c.thread != nil {
		<-c.thread
		c.thread = nil
	}

	logger.Info.Println("Camera client stopped")
}

// IsAlive checks if the camera client is running
func (c *PrinterCamera) IsAlive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.alive
}

// GetFrame gets the last camera frame as base64 encoded string
func (c *PrinterCamera) GetFrame() (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastFrame == nil {
		return "", fmt.Errorf("no frame available")
	}

	return base64.StdEncoding.EncodeToString(c.lastFrame), nil
}

// GetFrameBytes gets the last camera frame as bytes
func (c *PrinterCamera) GetFrameBytes() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastFrame == nil {
		return nil, fmt.Errorf("no frame available")
	}

	frameCopy := make([]byte, len(c.lastFrame))
	copy(frameCopy, c.lastFrame)
	return frameCopy, nil
}

// buildAuthData builds the authentication data for the camera connection
func (c *PrinterCamera) buildAuthData() []byte {
	authData := make([]byte, 0, 80)

	// Write header fields (little-endian)
	authData = append(authData, byte(0x40), 0x00, 0x00, 0x00) // 0x40
	authData = append(authData, byte(0x00), 0x30, 0x00, 0x00) // 0x3000
	authData = append(authData, byte(0x00), 0x00, 0x00, 0x00) // 0
	authData = append(authData, byte(0x00), 0x00, 0x00, 0x00) // 0

	// Write username (32 bytes, padded with zeros)
	usernameBytes := make([]byte, 32)
	copy(usernameBytes, []byte(c.username))
	authData = append(authData, usernameBytes...)

	// Write access code (32 bytes, padded with zeros)
	accessBytes := make([]byte, 32)
	copy(accessBytes, []byte(c.accessCode))
	authData = append(authData, accessBytes...)

	return authData
}

// retriever is the main camera retrieval loop
func (c *PrinterCamera) retriever() {
	defer close(c.thread)

	authData := c.buildAuthData()
	connectAttempts := 0

	jpegStart := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	jpegEnd := []byte{0xFF, 0xD9}

	readChunkSize := 4096

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	for c.alive {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.hostname, c.port))
		if err != nil {
			logger.Error.Printf("Error connecting to camera: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		tlsConn := tls.Client(conn, tlsConfig)
		err = tlsConn.Handshake()
		if err != nil {
			logger.Error.Printf("TLS handshake error: %v", err)
			tlsConn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		logger.Info.Println("Attempting to connect...")

		_, err = tlsConn.Write(authData)
		if err != nil {
			logger.Error.Printf("Error writing auth data: %v", err)
			tlsConn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		var img []byte
		var payloadSize int

		// Set read deadline
		tlsConn.SetReadDeadline(time.Now().Add(5 * time.Second))

		for c.alive {
			// Check if we should stop
			select {
			case <-c.stopChan:
				tlsConn.Close()
				return
			default:
			}

			buf := make([]byte, readChunkSize)
			n, err := tlsConn.Read(buf)

			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					time.Sleep(1 * time.Second)
					continue
				}
				if err == io.EOF {
					logger.Error.Println("Connection closed by server")
					break
				}
				logger.Error.Printf("Read error: %v", err)
				time.Sleep(1 * time.Second)
				break
			}

			logger.Debug.Printf("Read chunk: %d bytes", n)

			if img != nil && n > 0 {
				logger.Debug.Println("Appending to image")
				img = append(img, buf[:n]...)

				if len(img) > payloadSize {
					img = nil
				} else if len(img) == payloadSize {
					// Validate JPEG
					if len(img) >= 4 && string(img[:4]) != string(jpegStart) {
						img = nil
					} else if len(img) >= 2 && string(img[len(img)-2:]) != string(jpegEnd) {
						img = nil
					} else {
						c.mu.Lock()
						c.lastFrame = make([]byte, len(img))
						copy(c.lastFrame, img)
						c.mu.Unlock()
						img = nil
					}
				}
			} else if n == 16 {
				logger.Debug.Println("Got header")
				connectAttempts = 0
				img = make([]byte, 0)
				// Payload size is in bytes 0-3 (little-endian, 3 bytes)
				payloadSize = int(binary.LittleEndian.Uint32(append(buf[:3], 0)))
			} else if n == 0 {
				time.Sleep(5 * time.Second)
				logger.Error.Println("Wrong access code or IP")
				break
			} else {
				logger.Error.Println("Something bad happened")
				time.Sleep(1 * time.Second)
				break
			}
		}

		tlsConn.Close()

		if connectAttempts > 10 {
			logger.Error.Println("Too many connection attempts, reconnecting...")
			time.Sleep(5 * time.Second)
		}
	}
}
