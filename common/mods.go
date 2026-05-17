package common

type RcmfFile struct {
	Spec  string     `json:"spec"` //spec defines what RCMF spec to use. It should be a link to a document describing the spec. For this implementation, it should be set to "https://rcmf.axell.me/v1"
	Rules []RcmfRule `json:"rules"` //Rules describes the rules of this mod
}

type Sources struct { //Sources describes what a rule should modify or replace. Each rule may have multiple sources of multiple different types.
	Expressions []string `json:"expressions"` //This is reserved for future use. This can be ignored.
	Ids []int `json:"ids"` //Ids describes roblox asset ids to modify.
	Types []int `json:"types"` //Types describes asset types to modify in the form of asset type ids.
	Files []string `json:"files"` //Files describes files on the disk to modify. These strings can either be a path relative to the Roblox install directory (eg. version-xxxxxxxxxx), or it can be a special string that points to another path not in the roblox install directory. Special strings are "GlobalBasicSettings_13.xml" and "GlobalSettings_13.xml"
}

type RcmfRule struct { //Rules describe what to modify/replace and what should be changed
	Sources  Sources `json:"sources"` //Sources describes the sources of this rule
	Data RcmfData `json:"data"`//Data describes the data of this rule
}

type RcmfData struct { //Data describes what a rule should replace its Sources with or how to modify its sources. 
	Blob  []byte `json:"blob"` //Blob describes a binary object to replace the Sources with.
	Key   string `json:"key"` //Key describes a part of a xml or json document to modify. If the Sources are xml documents, key should be an XPath expression. If the Sources are json, a dot-seperated path should be used. Key should be empty when Blob is set.
	Value string `json:"value"` //Value describes the value to set a part of a xml or json document to when Key is set.
}

//This table describes what asset type each asset type id represent:
/*
Name	Value

Image	1	
TShirt	2	
Audio	3	
Mesh	4	
Lua	5	
Hat	8	
Place	9	
Model	10	
Shirt	11	
Pants	12	
Decal	13	
Head	17	
Face	18	
Gear	19	
Badge	21	
Animation	24	
Torso	27	
RightArm	28	
LeftArm	29	
LeftLeg	30	
RightLeg	31	
Package	32	
GamePass	34	
Plugin	38	
MeshPart	40	
HairAccessory	41	
FaceAccessory	42	
NeckAccessory	43	
ShoulderAccessory	44	
FrontAccessory	45	
BackAccessory	46	
WaistAccessory	47	
ClimbAnimation	48	
DeathAnimation	49	
FallAnimation	50	
IdleAnimation	51	
JumpAnimation	52	
RunAnimation	53	
SwimAnimation	54	
WalkAnimation	55	
PoseAnimation	56	
EarAccessory	57	
EyeAccessory	58	
EmoteAnimation	61	
Video	62	
TShirtAccessory	64	
ShirtAccessory	65	
PantsAccessory	66	
JacketAccessory	67	
SweaterAccessory	68	
ShortsAccessory	69	
LeftShoeAccessory	70	
RightShoeAccessory	71	
DressSkirtAccessory	72	
FontFamily	73	
EyebrowAccessory	76	
EyelashAccessory	77	
MoodAnimation	78	
DynamicHead	79	
FaceMakeup	88	
LipMakeup	89	
EyeMakeup	90	
VoxelFragment	91	
*/