package email

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type tcpServerMock struct {
	listen net.Listener
}

type tcpSession struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func newTCPServerMock(addr string) (*tcpServerMock, error) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	out := &tcpServerMock{
		listen: listen,
	}
	return out, nil
}

func (m *tcpServerMock) Start() {
	func() {
		for {
			conn, err := m.listen.Accept()
			sess := &tcpSession{
				conn:   conn,
				reader: bufio.NewReader(conn),
				writer: bufio.NewWriter(conn),
			}
			if err != nil {
				fmt.Printf("MOCK.server.loop.error: %+v\n", err)
				break
			}
			go func() {
				m.handle(sess)
				err := sess.conn.Close()
				if err != nil {
					fmt.Printf("MOCK.server.start.Close.error: %v\n", err)
				}
			}()
		}
	}()

}

func (sess *tcpSession) write(str string) {
	_, err := sess.writer.WriteString(fmt.Sprintf("%s\n", str))
	if err != nil {
		fmt.Printf("MOCK.server.sess.S error %+v\n", err)
	}
	err = sess.writer.Flush()
	if err != nil {
		fmt.Printf("MOCK.server.sess.S error %+v\n", err)
	}

	fmt.Printf("MOCK.server.sess.S -> %s\n", str)
}

func (*tcpServerMock) handle(sess *tcpSession) {
	defer sess.conn.Close()
	fmt.Printf("MOCK.server.sess.new: %+v\n", sess.conn.RemoteAddr())
	sess.write("220 localhost")
	for {
		var err error
		readLine, _, err := sess.reader.ReadLine()
		if err != nil {
			fmt.Printf("MOCK.server.error: %+v", err)
		}
		fmt.Printf("MOCK.server.sess.C -> [%s] \n", string(readLine))
		request := string(readLine)

		commandStr := strings.Split(strings.TrimSpace(request), " ")
		cmd := strings.ToUpper(commandStr[0])
		switch cmd {
		case "QUIT":
			sess.write("221 Bye")
			return
		case "HELO", "EHLO":
			sess.write("250-smtp.mock.com")
			sess.write("250-SIZE 1000000")
			sess.write("250 AUTH LOGIN PLAIN CRAM-MD5")
		case "AUTH":
			fmt.Printf("MOCK.server.sess.C.AUTH: %+v\n", commandStr[1:])
			if commandStr[2] == "AHVzZXIAcGFzcw==" {
				sess.write("235 2.7.0 Authentication successful")
			} else {
				sess.write("535 Authentication failed")
			}
		case "*":
			sess.write("535 Authentication failed")
		case "":
			sess.write("502 Command unrecognized. Available commands: HELO, EHLO, MAIL FROM:, RCPT TO:, DATA, RSET, NOOP, QUIT")
		}
	}

}

func (m *tcpServerMock) Close() {
	err := m.listen.Close()
	if err != nil {
		fmt.Printf("MOCK.server.error: %+v\n", err)
	}
}
