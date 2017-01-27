package douyu

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	log "qiniupkg.com/x/log.v7"
)

func OpenDanmu(rid int64) (dy *Douyu, err error) {
	dy = new(Douyu)
	dy.RoomID = rid
	dailer := &net.Dialer{
		Timeout: time.Second * 10,
	}
	if dy.Conn, err = dailer.Dial("tcp", "openbarrage.douyutv.com:8601"); err != nil {
		dy = nil
		return
	}
	dy.login()
	return
}

func (dy *Douyu) login() error {
	var resp [1024]byte
	reqData := append([]byte("type@=loginreq/roomid@="), number2bytes(dy.RoomID)...)

	reqData = PackRequest(reqData)
	dy.Write(reqData)
	cnt, err := dy.Read(resp[:])
	if err != nil {
		log.Error("login failed. err: ", err)
		return err
	}
	if cnt > 12 {
		dy.DouyuLoginResponse = parseLoginResponse(resp[12:cnt])
	} else {
		return errors.New("return nothing")
	}
	return nil
}

func parseLoginResponse(data []byte) *DouyuLoginResponse {
	fmt.Printf("data: %s", data)
	resp := new(DouyuLoginResponse)
	lines := bytes.Split(data, []byte{'/'})
	for _, line := range lines {
		if i := bytes.Index(line, []byte("@=")); i > 0 {
			k := line[:i]
			v := line[i+2:]
			switch string(k) {
			case "type":
				resp.Type = v
			case "userid":
				resp.UserID = bytes2number(v)
			case "roomgroup":
				resp.RoomGroup = bytes2number(v)
			case "pg":
				resp.Pg = bytes2number(v)
			case "sessionid":
				resp.SessionID = bytes2number(v)
			case "username":
				resp.Username = v
			case "nickname":
				resp.Nickname = v
			case "live_stat":
				resp.LiveStat = bytes2number(v) != 0
			case "is_illegal":
			case "ill_ct":
			case "is_signined":
				resp.IsSigned = bytes2number(v) != 0
			case "signin_count":
				resp.SignedCount = bytes2number(v)
			case "npv":
				resp.NeedPhoneVerify = bytes2number(v) != 0
			case "best_dlev":
				resp.BestDlev = bytes2number(v)
			case "cur_lev":
				resp.CurLev = bytes2number(v)
			case "nrc":
			case "sid":
			case "code":
				resp.ErrCode = bytes2number(v)
			default:
				fmt.Printf("%s", k)
			}
		}

	}
	return resp
}