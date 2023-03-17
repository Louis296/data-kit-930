package data_kit_930

// DataSet 文件数据集
type DataSet struct {
	PublicInfo      *PublicInfo
	DeviceInfo      *DeviceInfo
	AcquisitionInfo *AcquisitionInfo
	ImageInfo       *ImageInfo
	DataInfo        *DataInfo
}

// PublicInfo 1.2.1 公共信息
type PublicInfo struct {
	HeaderCRC uint16
	Length    uint32
	Type      uint16
}

// DeviceInfo 1.2.2 设备信息
type DeviceInfo struct {
	Length            uint32
	Device            string
	Serial            string
	AxisDetectors     uint16
	TransDetectors    uint16
	DetectorsRings    uint16
	DetectorsChannels uint16
	IPCounts          uint16
	IPStart           uint16
	ChannelCounts     uint16
	ChannelStart      uint16
	MVTThresholds     []float32
	MVTParameters     []float32
}

// AcquisitionInfo 1.2.3 采集信息
type AcquisitionInfo struct {
	Length             uint32
	Isotope            uint16
	Activity           float32
	InjectTime         string
	Time               string
	Duration           uint16
	TimeWindow         float32
	DelayWindow        float32
	XTalkWindow        float32
	EnergyWindow       []uint32
	PositionWindow     uint16
	Corrected          uint16
	TablePosition      float32
	TableHeight        float32
	PETCTSpacing       float32
	TableCount         uint16
	TableIndex         uint16
	ScanLengthPerTable float32
	PatientID          string
	StudyID            string
	PatientName        string
	PatientSex         string
	PatientHeight      float32
	PatientWeight      float32
}

// ImageInfo 1.2.4 图像信息
type ImageInfo struct {
	Length               uint32
	ImageSizeRows        uint16
	ImageSizeCols        uint16
	ImageSizeSlices      uint16
	ImageRowPixelSize    float32
	ImageColumnPixelSize float32
	ImageSliceThickness  float32
	ReconMethod          string
	MaxRingDiffNum       uint16
	SubsetNum            uint16
	IterNum              uint16
	AttnCalibration      uint16
	ScatCalibration      uint16
	ScatPara             []float32
	TVPara               []float32
	PetCtFovOffset       []float32
	CtRotationAngle      float32
	SeriesNumber         uint16
	ReconSoftwareVersion string
	PromptsCounts        uint32
	DelayCounts          uint32
}

type DataInfo struct {
	Length     uint32
	DataLength uint32
	CRC        uint16
}
