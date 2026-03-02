package beanstalkd

import (
	"net/url"

	"github.com/andeya/pholcus/config"
	"github.com/andeya/pholcus/logs"
	"github.com/kr/beanstalk"
	"github.com/pkg/errors"
)

// BeanstalkdClient wraps a beanstalk connection and tube for job queuing.
type BeanstalkdClient struct {
	Conn *beanstalk.Conn
	Tube string
}

// New creates a new BeanstalkdClient using config.BeanstalkdHost and config.BeanstalkdTube.
func New() (*BeanstalkdClient, error) {
	tmp := new(BeanstalkdClient)
	host := config.BeanstalkdHost
	if host == "" {
		return nil, errors.New("beanstalk host is empty")
	}
	tube := config.BeanstalkdTube
	if tube == "" {
		return nil, errors.New("tube name is empty")
	}
	conn, err := beanstalk.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	tmp.Tube = tube
	tmp.Conn = conn
	return tmp, nil
}

// Close closes the beanstalk connection.
func (srv *BeanstalkdClient) Close() {
	if srv.Conn != nil {
		srv.Conn.Close()
	}
}

// Send encodes content as URL values and puts it into the configured tube.
func (srv *BeanstalkdClient) Send(content url.Values) {
	if srv.Conn == nil {
		return
	}
	data := content.Encode()
	tube := &beanstalk.Tube{Conn: srv.Conn, Name: srv.Tube}

	_, err := tube.Put([]byte(data), 1, 0, 0)
	if err != nil {
		logs.Log.Error("beanstalkd write error: %v, content=%s", err, data)
		return
	}
}
