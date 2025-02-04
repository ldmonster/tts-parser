package module

type Module struct {
	SaveName      string `json:"SaveName"`
	EpochTime     int    `json:"EpochTime"`
	Date          string `json:"Date"`
	VersionNumber string `json:"VersionNumber"`

	// GameMode       string        `json:"GameMode"`
	// GameType       string        `json:"GameType"`
	// GameComplexity string        `json:"GameComplexity"`
	// PlayingTime    []int         `json:"PlayingTime"`
	// PlayerCounts   []int         `json:"PlayerCounts"`
	// Tags           []string      `json:"Tags"`
	// Gravity        float64       `json:"Gravity"`
	// PlayArea       float64       `json:"PlayArea"`
	// Table          string        `json:"Table"`
	// Sky            string        `json:"Sky"`
	// Note           string        `json:"Note"`
	// TabStates      TabStates     `json:"TabStates"`
	// Grid           Grid          `json:"Grid"`
	// Lighting       Lighting      `json:"Lighting"`
	// Hands          Hands         `json:"Hands"`
	// ComponentTags  ComponentTags `json:"ComponentTags"`
	// Turns          Turns         `json:"Turns"`
	// DecalPallet    []any         `json:"DecalPallet"`
	// LuaScript      string        `json:"LuaScript"`
	// LuaScriptState string        `json:"LuaScriptState"`
	// XMLUI          string        `json:"XmlUI"`

	TableURL       string         `json:"TableURL"`
	SkyURL         string         `json:"SkyURL"`
	CustomUIAssets CustomUIAssets `json:"CustomUIAssets,omitempty"`
	MusicPlayer    *MusicPlayer   `json:"MusicPlayer"`

	Objects []Object `json:"ObjectStates"`
}

type Object struct {
	// GUID                 string               `json:"GUID"`
	// Name                 string               `json:"Name"`
	// Transform            Transform            `json:"Transform"`
	// Nickname             string               `json:"Nickname"`
	// Description          string               `json:"Description"`
	// GMNotes              string               `json:"GMNotes"`
	// AltLookAngle         AltLookAngle         `json:"AltLookAngle"`
	// ColorDiffuse         Color                `json:"ColorDiffuse,omitempty"`
	// LayoutGroupSortIndex int                  `json:"LayoutGroupSortIndex"`
	// Value                int                  `json:"Value"`
	// Locked               bool                 `json:"Locked"`
	// Grid                 bool                 `json:"Grid"`
	// Snap                 bool                 `json:"Snap"`
	// IgnoreFoW            bool                 `json:"IgnoreFoW"`
	// MeasureMovement      bool                 `json:"MeasureMovement"`
	// DragSelectable       bool                 `json:"DragSelectable"`
	// Autoraise            bool                 `json:"Autoraise"`
	// Sticky               bool                 `json:"Sticky"`
	// Tooltip              bool                 `json:"Tooltip"`
	// GridProjection       bool                 `json:"GridProjection"`
	// HideWhenFaceDown     bool                 `json:"HideWhenFaceDown"`
	// Hands                bool                 `json:"Hands"`
	// FogColor             string               `json:"FogColor,omitempty"`
	// LuaScript            string               `json:"LuaScript"`
	// LuaScriptState       string               `json:"LuaScriptState"`
	// XMLUI                string               `json:"XmlUI"`
	// PhysicsMaterial      PhysicsMaterial      `json:"PhysicsMaterial,omitempty"`
	// Rigidbody            Rigidbody            `json:"Rigidbody,omitempty"`
	// MaterialIndex        int                  `json:"MaterialIndex,omitempty"`
	// MeshIndex            int                  `json:"MeshIndex,omitempty"`
	// Bag                  Bag                  `json:"Bag,omitempty"`
	// RotationValues       []RotationValues     `json:"RotationValues,omitempty"`
	// Clock                Clock                `json:"Clock,omitempty"`
	// SidewaysCard         bool                 `json:"SidewaysCard,omitempty"`
	// DeckIDs              []int                `json:"DeckIDs,omitempty"`
	// CardID               int                  `json:"CardID,omitempty"`
	// Number               int                  `json:"Number,omitempty"`
	// Text                 Text                 `json:"Text,omitempty"`
	// AttachedSnapPoints   []AttachedSnapPoints `json:"AttachedSnapPoints,omitempty"`

	CustomDeck        CustomDeck         `json:"CustomDeck,omitempty"`
	AttachedDecals    AttachedDecals     `json:"AttachedDecals,omitempty"`
	CustomAssetbundle *CustomAssetbundle `json:"CustomAssetbundle,omitempty"`
	CustomUIAssets    CustomUIAssets     `json:"CustomUIAssets,omitempty"`
	CustomMesh        *CustomMesh        `json:"CustomMesh,omitempty"`
	CustomImage       *CustomImage       `json:"CustomImage,omitempty"`
	CustomPDF         *CustomPDF         `json:"CustomPDF,omitempty"`

	States           States   `json:"States,omitempty"`
	ContainedObjects []Object `json:"ContainedObjects,omitempty"`
	ChildObjects     []Object `json:"ChildObjects,omitempty"`
}

// type Transform struct {
// 	PosX float64 `json:"posX"`
// 	PosY float64 `json:"posY"`
// 	PosZ float64 `json:"posZ"`
// 	RotX float64 `json:"rotX"`
// 	RotY float64 `json:"rotY"`
// 	RotZ float64 `json:"rotZ"`
// 	ScaleX float64 `json:"scaleX"`
// 	ScaleY float64 `json:"scaleY"`
// 	ScaleZ float64 `json:"scaleZ"`
// }

// type AltLookAngle struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// 	Z float64 `json:"z"`
// }

type CustomUIAssets []CustomUIAsset

type CustomUIAsset struct {
	Type int    `json:"Type"`
	Name string `json:"Name"`
	URL  string `json:"URL"`
}

type CustomAssetbundle struct {
	AssetbundleURL          string `json:"AssetbundleURL"`
	AssetbundleSecondaryURL string `json:"AssetbundleSecondaryURL"`
	MaterialIndex           int    `json:"MaterialIndex"`
	TypeIndex               int    `json:"TypeIndex"`
	LoopingEffectIndex      int    `json:"LoopingEffectIndex"`
}

// type PhysicsMaterial struct {
// 	StaticFriction  float64 `json:"StaticFriction"`
// 	DynamicFriction float64 `json:"DynamicFriction"`
// 	Bounciness      float64 `json:"Bounciness"`
// 	FrictionCombine int     `json:"FrictionCombine"`
// 	BounceCombine   int     `json:"BounceCombine"`
// }

// type Rigidbody struct {
// 	Mass        float64 `json:"Mass"`
// 	Drag        float64 `json:"Drag"`
// 	AngularDrag float64 `json:"AngularDrag"`
// 	UseGravity  bool    `json:"UseGravity"`
// }

type States map[string]Object

type CustomMesh struct {
	MeshURL     string `json:"MeshURL"`
	DiffuseURL  string `json:"DiffuseURL"`
	NormalURL   string `json:"NormalURL"`
	ColliderURL string `json:"ColliderURL"`
	// Convex        bool         `json:"Convex"`
	// MaterialIndex int          `json:"MaterialIndex"`
	// TypeIndex     int          `json:"TypeIndex"`
	// CustomShader  CustomShader `json:"CustomShader"`
	// CastShadows   bool         `json:"CastShadows"`
}

// type Bag struct {
// 	Order int `json:"Order"`
// }

// type CustomShader struct {
// 	SpecularColor     Color   `json:"SpecularColor"`
// 	SpecularIntensity float64 `json:"SpecularIntensity"`
// 	SpecularSharpness float64 `json:"SpecularSharpness"`
// 	FresnelStrength   float64 `json:"FresnelStrength"`
// }

type CustomImage struct {
	ImageURL          string `json:"ImageURL"`
	ImageSecondaryURL string `json:"ImageSecondaryURL"`
	// ImageScalar       float64    `json:"ImageScalar"`
	// WidthScale        float64    `json:"WidthScale"`
	// CustomDice        CustomDice `json:"CustomDice"`
}

// type CustomDice struct {
// 	Type int `json:"Type"`
// }

// type Rotation struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// 	Z float64 `json:"z"`
// }

// type RotationValues struct {
// 	Value    string   `json:"Value"`
// 	Rotation Rotation `json:"Rotation"`
// }

// type CustomTile struct {
// 	Type      int     `json:"Type"`
// 	Thickness float64 `json:"Thickness"`
// 	Stackable bool    `json:"Stackable"`
// 	Stretch   bool    `json:"Stretch"`
// }

// type Clock struct {
// 	Mode          int  `json:"Mode"`
// 	SecondsPassed int  `json:"SecondsPassed"`
// 	Paused        bool `json:"Paused"`
// }

type CustomCard struct {
	FaceURL string `json:"FaceURL"`
	BackURL string `json:"BackURL"`
	// NumWidth     int    `json:"NumWidth"`
	// NumHeight    int    `json:"NumHeight"`
	// BackIsHidden bool   `json:"BackIsHidden"`
	// UniqueBack   bool   `json:"UniqueBack"`
	// Type         int    `json:"Type"`
}

type CustomDeck map[string]CustomCard

type AttachedDecals []Decal

type Decal struct {
	CustomDecal *CustomDecal `json:"CustomDecal,omitempty"`
}

type CustomDecal struct {
	Name     string `json:"Name"`
	ImageURL string `json:"ImageURL"`
	// Size     float64 `json:"Size"`
}

// type Text struct {
// 	Text       string `json:"Text"`
// 	Colorstate Color  `json:"colorstate"`
// 	FontSize   int    `json:"fontSize"`
// }

// type Position struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// 	Z float64 `json:"z"`
// }

// type AttachedSnapPoints struct {
// 	Position Position `json:"Position"`
// }

type CustomPDF struct {
	PDFURL string `json:"PDFUrl"`
	// PDFPassword   string `json:"PDFPassword"`
	// PDFPage       int    `json:"PDFPage"`
	// PDFPageOffset int    `json:"PDFPageOffset"`
}

// type TabState struct {
// 	Title        string `json:"title"`
// 	Body         string `json:"body"`
// 	Color        string `json:"color"`
// 	VisibleColor Color  `json:"visibleColor"`
// 	ID           int    `json:"id"`
// }

// type TabStates map[string]TabState

type AudioLibrary map[string]string

type MusicPlayer struct {
	// RepeatSong        bool           `json:"RepeatSong"`
	// PlaylistEntry     int            `json:"PlaylistEntry"`
	// CurrentAudioTitle string         `json:"CurrentAudioTitle"`
	// CurrentAudioURL   string         `json:"CurrentAudioURL"`
	AudioLibrary []AudioLibrary `json:"AudioLibrary"`
}

// type Color struct {
// 	R float64 `json:"r"`
// 	G float64 `json:"g"`
// 	B float64 `json:"b"`
// 	A float64 `json:"a,omitempty"`
// }

type PosOffset struct {
	// X float64 `json:"x"`
	// Y float64 `json:"y"`
	// Z float64 `json:"z"`
}

// type Grid struct {
// 	Type         int       `json:"Type"`
// 	Lines        bool      `json:"Lines"`
// 	Color        Color     `json:"Color"`
// 	Opacity      float64   `json:"Opacity"`
// 	ThickLines   bool      `json:"ThickLines"`
// 	Snapping     bool      `json:"Snapping"`
// 	Offset       bool      `json:"Offset"`
// 	BothSnapping bool      `json:"BothSnapping"`
// 	XSize        float64   `json:"xSize"`
// 	YSize        float64   `json:"ySize"`
// 	PosOffset    PosOffset `json:"PosOffset"`
// }

// type Lighting struct {
// 	LightIntensity      float64 `json:"LightIntensity"`
// 	LightColor          Color   `json:"LightColor"`
// 	AmbientIntensity    float64 `json:"AmbientIntensity"`
// 	AmbientType         int     `json:"AmbientType"`
// 	AmbientSkyColor     Color   `json:"AmbientSkyColor"`
// 	AmbientEquatorColor Color   `json:"AmbientEquatorColor"`
// 	AmbientGroundColor  Color   `json:"AmbientGroundColor"`
// 	ReflectionIntensity float64 `json:"ReflectionIntensity"`
// 	LutIndex            int     `json:"LutIndex"`
// 	LutContribution     float64 `json:"LutContribution"`
// 	LutURL              string  `json:"LutURL"`
// }

// type Hands struct {
// 	Enable        bool `json:"Enable"`
// 	DisableUnused bool `json:"DisableUnused"`
// 	Hiding        int  `json:"Hiding"`
// }

// type ComponentTags struct {
// 	Labels []any `json:"labels"`
// }

// type Turns struct {
// 	Enable              bool   `json:"Enable"`
// 	Type                int    `json:"Type"`
// 	TurnOrder           []any  `json:"TurnOrder"`
// 	Reverse             bool   `json:"Reverse"`
// 	SkipEmpty           bool   `json:"SkipEmpty"`
// 	DisableInteractions bool   `json:"DisableInteractions"`
// 	PassTurns           bool   `json:"PassTurns"`
// 	TurnColor           string `json:"TurnColor"`
// }
