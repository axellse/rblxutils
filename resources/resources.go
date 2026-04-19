package resources

import _ "embed"

//go:embed pic_1.png
var CatPic1 []byte
//go:embed pic_2.png
var CatPic2 []byte
//go:embed pic_3.png
var CatPic3 []byte
//go:embed pic_4.png
var CatPic4 []byte
//go:embed pic_5.png
var CatPic5 []byte
//go:embed pic_6.png
var CatPic6 []byte
//go:embed pic_7.png
var CatPic7 []byte
//go:embed packagemap.json
var PackageMap []byte
//go:embed AppSettings.xml
var AppSettings []byte

//go:embed pic_4.png
var ProgramLogo []byte
//go:embed welcome_cat.png
var WelcomeCatImage []byte
var ProgramName = "rblxutils"