package resources

import _ "embed"

//go:embed images/pic_1.png
var CatPic1 []byte
//go:embed images/pic_2.png
var CatPic2 []byte
//go:embed images/pic_3.png
var CatPic3 []byte
//go:embed images/pic_4.png
var CatPic4 []byte
//go:embed images/pic_5.png
var CatPic5 []byte
//go:embed images/pic_6.png
var CatPic6 []byte
//go:embed images/pic_7.png
var CatPic7 []byte

//go:embed packagemap.json
var PackageMap []byte
//go:embed AppSettings.xml
var AppSettings []byte

//go:embed images/logo.png
var ProgramLogo []byte
//go:embed images/logo.ico
var ProgramLogoIco []byte
//go:embed images/roblox_r.png
var RobloxRLogo []byte
//go:embed images/builder_club.png
var BuilderClubLogo []byte
//go:embed images/welcome_cat.png
var WelcomeCatImage []byte
//go:embed images/apartments_ljms.png
var ApartmentsLjmsImage []byte

//go:embed cryptography/update_key.pem
var UpdatePublicKey []byte
//go:embed cryptography/ca.crt
var CACert []byte

//go:embed cryptography/assetdelivery.roblox.com.crt
var AssetdeliveryCert []byte
//go:embed cryptography/assetdelivery.roblox.com.key
var AssetdeliveryKey []byte

//go:embed cryptography/fts.rbxcdn.com.crt
var RbxcdnCert []byte
//go:embed cryptography/fts.rbxcdn.com.key
var RbxcdnKey []byte

//go:embed update.bat
var UpdateBatch []byte
//go:embed version
var Version string