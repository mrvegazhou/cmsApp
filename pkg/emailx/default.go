package emailx

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
	"sync"
	"time"
)

func SendEmailByPool(title, content string, sendEmailers []string) error {
	// 创建有len(sendEmailers)个缓冲的通道
	nums := len(sendEmailers)
	ch := make(chan *email.Email, nums)
	emailConf := configs.App.Email
	// 连接池
	p, err := email.NewPool(
		fmt.Sprintf("%s:%s", emailConf.Smtp, emailConf.SmtpPort),
		emailConf.PoolSize,
		smtp.PlainAuth("", emailConf.Value, emailConf.Password, emailConf.Smtp),
	)
	if err != nil {
		errStr := fmt.Sprintf("%s %s", constant.SEND_EMAIL_POOL_ERR, err)
		return errors.New(errStr)
	}
	var wg sync.WaitGroup
	wg.Add(emailConf.PoolSize)
	for i := 0; i < emailConf.PoolSize; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
				}
			}
		}()
	}

	for i := 0; i < nums; i++ {
		e := email.NewEmail()
		e.From = fmt.Sprintf("%s <%s>", configs.App.Email.EmailName, configs.App.Email.Value)
		e.To = []string{sendEmailers[i]}
		e.Subject = title
		e.Text = []byte(content)
		ch <- e
	}

	close(ch)
	wg.Wait()
	return nil
}

func SendEmail(title, content, snedEmailer string) error {
	emailConf := configs.App.Email
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = fmt.Sprintf("%s <%s>", emailConf.EmailName, emailConf.Value)
	// 设置接收方的邮箱
	e.To = []string{snedEmailer}
	//设置主题
	e.Subject = title
	//设置文件发送的内容
	e.HTML = []byte(content)
	//设置服务器相关的配置
	err := e.Send(fmt.Sprintf("%s:%s", emailConf.Smtp, emailConf.SmtpPort), smtp.PlainAuth("", emailConf.Value, emailConf.Password, emailConf.Smtp))
	if err != nil {
		return err
	}
	return nil
}
