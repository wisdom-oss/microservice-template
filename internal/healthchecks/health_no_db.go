//go:build no_db

package healthchecks

// Base is a very basic healthcheck that pings the databse server and returns an
// error if the connection could not be established.
func Base() error {
	return nil
}
