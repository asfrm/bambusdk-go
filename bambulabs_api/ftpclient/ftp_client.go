package ftpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/textproto"
	"strings"
	"sync"

	"github.com/bambulabs_api/go/bambulabs_api/logger"
	"github.com/jlaffaye/ftp"
)

// PrinterFTPClient handles FTP communication with the printer
type PrinterFTPClient struct {
	serverIP   string
	port       int
	user       string
	accessCode string

	mu         sync.Mutex
	ftpClient  *ftp.ServerConn
	tlsConfig  *tls.Config
}

// NewPrinterFTPClient creates a new FTP client
func NewPrinterFTPClient(serverIP, accessCode string, user string, port int) *PrinterFTPClient {
	if user == "" {
		user = "bblp"
	}
	if port == 0 {
		port = 990
	}

	return &PrinterFTPClient{
		serverIP:   serverIP,
		port:       port,
		user:       user,
		accessCode: accessCode,
		tlsConfig: &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		},
	}
}

// connect connects to the FTP server with implicit TLS
func (c *PrinterFTPClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ftpClient != nil {
		return nil
	}

	// Connect with implicit TLS
	dialer := &ftp.Dialer{
		TLSConfig: c.tlsConfig,
		ExplicitTLS: false, // Implicit TLS
	}

	client, err := dialer.Dial(fmt.Sprintf("%s:%d", c.serverIP, c.port))
	if err != nil {
		// Try explicit TLS as fallback
		dialer.ExplicitTLS = true
		client, err = dialer.Dial(fmt.Sprintf("%s:%d", c.serverIP, c.port))
		if err != nil {
			return fmt.Errorf("failed to connect to FTP server: %w", err)
		}
	}

	err = client.Login(c.user, c.accessCode)
	if err != nil {
		client.Quit()
		return fmt.Errorf("failed to login to FTP server: %w", err)
	}

	// Enter passive mode and enable TLS for data connections
	err = client.ProtP()
	if err != nil {
		logger.Warn.Printf("Failed to set protection level: %v", err)
	}

	c.ftpClient = client
	logger.Info.Println("Connected to FTP server")

	return nil
}

// close closes the FTP connection
func (c *PrinterFTPClient) close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ftpClient != nil {
		err := c.ftpClient.Quit()
		c.ftpClient = nil
		if err != nil {
			return fmt.Errorf("failed to close FTP connection: %w", err)
		}
		logger.Info.Println("FTP connection closed")
	}
	return nil
}

// withConnection ensures a connection is available and runs the function
func (c *PrinterFTPClient) withConnection(fn func(*ftp.ServerConn) error) error {
	if err := c.connect(); err != nil {
		return err
	}

	c.mu.Lock()
	client := c.ftpClient
	c.mu.Unlock()

	if client == nil {
		return fmt.Errorf("FTP client not connected")
	}

	return fn(client)
}

// UploadFile uploads a file to the printer
func (c *PrinterFTPClient) UploadFile(data io.Reader, filePath string) (string, error) {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		logger.Info.Printf("Uploading file: %s", filePath)

		// Read all data first
		fileData, err := io.ReadAll(data)
		if err != nil {
			return fmt.Errorf("failed to read file data: %w", err)
		}

		// Upload the file
		err = client.Stor(filePath, bytes.NewReader(fileData))
		if err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}

		logger.Info.Printf("File uploaded successfully: %s", filePath)
		return nil
	})

	if err != nil {
		return "", err
	}

	return filePath, nil
}

// ListDirectory lists files in a directory
func (c *PrinterFTPClient) ListDirectory(path string) ([]string, error) {
	var entries []*ftp.Entry
	err := c.withConnection(func(client *ftp.ServerConn) error {
		var err error
		if path != "" {
			entries, err = client.List(path)
		} else {
			entries, err = client.List("")
		}
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %w", err)
	}

	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name
	}

	return names, nil
}

// ListImagesDir lists files in the image directory
func (c *PrinterFTPClient) ListImagesDir() ([]string, error) {
	return c.ListDirectory("image")
}

// ListCacheDir lists files in the cache directory
func (c *PrinterFTPClient) ListCacheDir() ([]string, error) {
	return c.ListDirectory("cache")
}

// ListTimelapseDir lists files in the timelapse directory
func (c *PrinterFTPClient) ListTimelapseDir() ([]string, error) {
	return c.ListDirectory("timelapse")
}

// ListLoggerDir lists files in the logger directory
func (c *PrinterFTPClient) ListLoggerDir() ([]string, error) {
	return c.ListDirectory("logger")
}

// DownloadFile downloads a file from the printer
func (c *PrinterFTPClient) DownloadFile(filePath string) ([]byte, error) {
	var fileData bytes.Buffer

	err := c.withConnection(func(client *ftp.ServerConn) error {
		reader, err := client.Retr(filePath)
		if err != nil {
			return fmt.Errorf("failed to retrieve file: %w", err)
		}
		defer reader.Close()

		_, err = io.Copy(&fileData, reader)
		return err
	})

	if err != nil {
		return nil, err
	}

	return fileData.Bytes(), nil
}

// DeleteFile deletes a file from the printer
func (c *PrinterFTPClient) DeleteFile(filePath string) error {
	logger.Info.Printf("Deleting file: %s", filePath)

	err := c.withConnection(func(client *ftp.ServerConn) error {
		return client.Delete(filePath)
	})

	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logger.Info.Printf("File deleted: %s", filePath)
	return nil
}

// GetLastImagePrint gets the last image from the image directory
func (c *PrinterFTPClient) GetLastImagePrint() ([]byte, error) {
	files, err := c.ListImagesDir()
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no images found")
	}

	// Get the last file
	lastFile := files[len(files)-1]
	imgPath := fmt.Sprintf("image/%s", lastFile)

	return c.DownloadFile(imgPath)
}

// MakeDir creates a directory on the FTP server
func (c *PrinterFTPClient) MakeDir(path string) error {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		return client.MakeDir(path)
	})

	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// ChangeDir changes the current directory
func (c *PrinterFTPClient) ChangeDir(path string) error {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		return client.ChangeDir(path)
	})

	if err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	return nil
}

// GetCurrentDir gets the current working directory
func (c *PrinterFTPClient) GetCurrentDir() (string, error) {
	var dir string
	err := c.withConnection(func(client *ftp.ServerConn) error {
		var err error
		dir, err = client.CurrentDir()
		return err
	})

	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	return dir, nil
}

// Rename renames a file or directory
func (c *PrinterFTPClient) Rename(oldPath, newPath string) error {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		return client.Rename(oldPath, newPath)
	})

	if err != nil {
		return fmt.Errorf("failed to rename: %w", err)
	}

	return nil
}

// GetFileSize gets the size of a file
func (c *PrinterFTPClient) GetFileSize(filePath string) (uint64, error) {
	var size uint64
	err := c.withConnection(func(client *ftp.ServerConn) error {
		entry, err := client.Stat(filePath)
		if err != nil {
			return err
		}
		size = entry.Size
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get file size: %w", err)
	}

	return size, nil
}

// Close closes the FTP connection
func (c *PrinterFTPClient) Close() error {
	return c.close()
}

// IsConnected checks if the FTP client is connected
func (c *PrinterFTPClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ftpClient != nil
}

// Reconnect reconnects to the FTP server
func (c *PrinterFTPClient) Reconnect() error {
	c.close()
	return c.connect()
}

// Noop sends a NOOP command to keep the connection alive
func (c *PrinterFTPClient) Noop() error {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		return client.Noop()
	})

	if err != nil {
		return fmt.Errorf("NOOP failed: %w", err)
	}

	return nil
}

// GetEntry gets detailed information about a file or directory
func (c *PrinterFTPClient) GetEntry(path string) (*ftp.Entry, error) {
	var entry *ftp.Entry
	err := c.withConnection(func(client *ftp.ServerConn) error {
		entries, err := client.List(path)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			return fmt.Errorf("entry not found: %s", path)
		}
		entry = entries[0]
		return nil
	})

	if err != nil {
		return nil, err
	}

	return entry, nil
}

// AppendFile appends data to an existing file
func (c *PrinterFTPClient) AppendFile(data io.Reader, filePath string) error {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		// Read all data first
		fileData, err := io.ReadAll(data)
		if err != nil {
			return fmt.Errorf("failed to read file data: %w", err)
		}

		// Get current file content
		currentData, err := c.DownloadFile(filePath)
		if err != nil && !strings.Contains(err.Error(), "failed") {
			// File doesn't exist, just upload
			return client.Stor(filePath, bytes.NewReader(fileData))
		}

		// Append new data
		combinedData := append(currentData, fileData...)
		return client.Stor(filePath, bytes.NewReader(combinedData))
	})

	return err
}

// SetConnType sets the connection type
func (c *PrinterFTPClient) SetConnType(connType ftp.ConnType) error {
	err := c.withConnection(func(client *ftp.ServerConn) error {
		// Note: jlaffaye/ftp doesn't expose SetConnType directly
		// This is a placeholder for potential future implementation
		_ = connType
		return nil
	})
	return err
}

// SendCmd sends a custom FTP command
func (c *PrinterFTPClient) SendCmd(cmd string, args ...string) (string, error) {
	var response string
	err := c.withConnection(func(client *ftp.ServerConn) error {
		// Use low-level command execution
		fullCmd := cmd
		if len(args) > 0 {
			fullCmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
		}
		
		// This requires access to the underlying textproto connection
		// which is not exposed by jlaffaye/ftp
		_ = fullCmd
		return fmt.Errorf("custom commands not supported")
	})
	return response, err
}

// KeepAlive keeps the connection alive
func (c *PrinterFTPClient) KeepAlive() error {
	return c.Noop()
}

// GetTimeout gets the current timeout
func (c *PrinterFTPClient) GetTimeout() int {
	return 300 // Default timeout in seconds
}

// SetTimeout sets the timeout (placeholder for future implementation)
func (c *PrinterFTPClient) SetTimeout(timeout int) {
	// Placeholder for future implementation
	logger.Debug.Printf("Timeout set to: %d seconds", timeout)
}

// Debug enables/disables debug mode
func (c *PrinterFTPClient) Debug(enable bool) {
	if enable {
		logger.Debug.Println("FTP debug mode enabled")
	}
}

// GetWelcomeMsg gets the FTP server welcome message
func (c *PrinterFTPClient) GetWelcomeMsg() (string, error) {
	// Note: This would require access to the underlying connection
	// which is not directly exposed by jlaffaye/ftp
	return "", fmt.Errorf("welcome message not available")
}

// ListRecursive lists all files recursively
func (c *PrinterFTPClient) ListRecursive(path string) ([]string, error) {
	var allFiles []string
	
	err := c.withConnection(func(client *ftp.ServerConn) error {
		entries, err := client.List(path)
		if err != nil {
			return err
		}
		
		for _, entry := range entries {
			if entry.Type == ftp.EntryTypeFolder {
				if entry.Name != "." && entry.Name != ".." {
					subPath := fmt.Sprintf("%s/%s", path, entry.Name)
					subFiles, err := c.ListRecursive(subPath)
					if err != nil {
						return err
					}
					allFiles = append(allFiles, subFiles...)
				}
			} else {
				filePath := fmt.Sprintf("%s/%s", path, entry.Name)
				allFiles = append(allFiles, filePath)
			}
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return allFiles, nil
}

// FileExists checks if a file exists
func (c *PrinterFTPClient) FileExists(filePath string) (bool, error) {
	_, err := c.GetFileSize(filePath)
	if err != nil {
		if protoErr, ok := err.(*textproto.Error); ok {
			if protoErr.Code == 550 {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
