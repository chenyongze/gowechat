package gowechat

import (
	"testing"

	"github.com/astaxie/beego"
	"github.com/chenyongze/gowechat/wxcontext"
)

func TestGetQrcode(t *testing.T) {
	config := wxcontext.Config{
		AppID:     "wxb230e8e6858ad365",
		AppSecret: "587b86fa560021de690bbbb10e5e4afe",
		Token:     "xgk_test",
	}
	wc := NewWechat(config)
	beego.Debug("wechat's cache:", wc.Context.Cache)
	mp, _ := wc.MpMgr()
	qrcode, err := mp.GetQrcode().CreatePermanentQRCodeWithSceneString("aaa")
	beego.Debug("qrcode:", qrcode.ImageURL())
	beego.Debug("err:", err)
}
