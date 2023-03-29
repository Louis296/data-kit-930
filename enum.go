package data_kit_930

// 数据文件文件类型
const (
	RawDataType uint16 = iota
	ListmodeDataType
	MichDataType
	EnergyCalibrationMap
	TimeCalibrationMap
	EnergySpectrumData
)

// IP前缀
const ipPrefix = "192.168."
