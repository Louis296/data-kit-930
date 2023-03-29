package data_kit_930

import (
	"encoding/binary"
	"io"
	"os"
)

func ParseFile(path string) (*DataSet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	p := &Parser{
		reader:    file,
		byteOrder: binary.BigEndian,
		modifyStr: true,
	}
	return p.parser(), nil
}

type Parser struct {
	reader    io.Reader
	byteOrder binary.ByteOrder

	// 是否移除string末尾的空字符
	modifyStr bool
}

func (p *Parser) parser() *DataSet {
	publicInfo := p.parsePublicInfo()
	deviceInfo := p.parseDeviceInfo()
	var acquisitionInfo *AcquisitionInfo
	var imageInfo *ImageInfo
	var dataInfo *DataInfo
	switch publicInfo.Type {
	case RawDataType:
		acquisitionInfo = p.parseAcquisitionInfo()
		dataInfo = p.parseDataInfo()
	case ListmodeDataType:
		acquisitionInfo = p.parseAcquisitionInfo()
		dataInfo = p.parseDataInfo()
	case MichDataType:
		acquisitionInfo = p.parseAcquisitionInfo()
		dataInfo = p.parseDataInfo()
	case EnergyCalibrationMap:
		dataInfo = p.parseDataInfo()
	case TimeCalibrationMap:
		dataInfo = p.parseDataInfo()
	case EnergySpectrumData:
		dataInfo = p.parseDataInfo()
	default:
		acquisitionInfo = p.parseAcquisitionInfo()
		imageInfo = p.parseImageInfo()
		dataInfo = p.parseDataInfo()
	}
	return &DataSet{
		PublicInfo:      publicInfo,
		DeviceInfo:      deviceInfo,
		AcquisitionInfo: acquisitionInfo,
		ImageInfo:       imageInfo,
		DataInfo:        dataInfo,
	}
}

func (p *Parser) parsePublicInfo() *PublicInfo {
	// skip magic keys
	_, _ = p.nextString(16)
	return &PublicInfo{
		HeaderCRC:       p.mustNextUint16(),
		Length:          p.mustNextUint32(),
		Type:            p.mustNextUint16(),
		SoftwareVersion: p.mustNextString(16),
		HeaderLength:    p.mustNextUint32(),
	}
}

func (p *Parser) parseDeviceInfo() *DeviceInfo {
	return &DeviceInfo{
		Length:            p.mustNextUint32(),
		Device:            p.mustNextString(16),
		Serial:            p.mustNextString(16),
		AxisDetectors:     p.mustNextUint16(),
		TransDetectors:    p.mustNextUint16(),
		DetectorsRings:    p.mustNextUint16(),
		DetectorsChannels: p.mustNextUint16(),
		IPCounts:          p.mustNextUint16(),
		IPStart:           p.mustNextUint16(),
		ChannelCounts:     p.mustNextUint16(),
		ChannelStart:      p.mustNextUint16(),
		MVTThresholds:     p.mustNextFloat32Slice(8),
		MVTParameters:     p.mustNextFloat32Slice(3),
	}
}

func (p *Parser) parseAcquisitionInfo() *AcquisitionInfo {
	return &AcquisitionInfo{
		Length:             p.mustNextUint32(),
		Isotope:            p.mustNextUint16(),
		Activity:           p.mustNextFloat32(),
		InjectTime:         p.mustNextString(16),
		Time:               p.mustNextString(16),
		Duration:           p.mustNextUint16(),
		TimeWindow:         p.mustNextFloat32(),
		DelayWindow:        p.mustNextFloat32(),
		XTalkWindow:        p.mustNextFloat32(),
		EnergyWindow:       []uint32{p.mustNextUint32(), p.mustNextUint32()},
		PositionWindow:     p.mustNextUint16(),
		Corrected:          p.mustNextUint16(),
		TablePosition:      p.mustNextFloat32(),
		TableHeight:        p.mustNextFloat32(),
		PETCTSpacing:       p.mustNextFloat32(),
		TableCount:         p.mustNextUint16(),
		TableIndex:         p.mustNextUint16(),
		ScanLengthPerTable: p.mustNextFloat32(),
		PatientID:          p.mustNextString(64),
		StudyID:            p.mustNextString(64),
		PatientName:        p.mustNextString(128),
		PatientSex:         p.mustNextString(8),
		PatientHeight:      p.mustNextFloat32(),
		PatientWeight:      p.mustNextFloat32(),
	}
}

func (p *Parser) parseImageInfo() *ImageInfo {
	return &ImageInfo{
		Length:               p.mustNextUint32(),
		ImageSizeRows:        p.mustNextUint16(),
		ImageSizeCols:        p.mustNextUint16(),
		ImageSizeSlices:      p.mustNextUint16(),
		ImageRowPixelSize:    p.mustNextFloat32(),
		ImageColumnPixelSize: p.mustNextFloat32(),
		ImageSliceThickness:  p.mustNextFloat32(),
		ReconMethod:          p.mustNextString(16),
		MaxRingDiffNum:       p.mustNextUint16(),
		SubsetNum:            p.mustNextUint16(),
		IterNum:              p.mustNextUint16(),
		AttnCalibration:      p.mustNextUint16(),
		ScatCalibration:      p.mustNextUint16(),
		ScatPara:             p.mustNextFloat32Slice(6),
		TVPara:               p.mustNextFloat32Slice(2),
		PetCtFovOffset:       p.mustNextFloat32Slice(3),
		CtRotationAngle:      p.mustNextFloat32(),
		SeriesNumber:         p.mustNextUint16(),
		ReconSoftwareVersion: p.mustNextString(16),
		PromptsCounts:        p.mustNextUint32(),
		DelayCounts:          p.mustNextUint32(),
	}
}

func (p *Parser) parseDataInfo() *DataInfo {
	return &DataInfo{
		Length:     p.mustNextUint32(),
		DataLength: p.mustNextUint32(),
		CRC:        p.mustNextUint16(),
	}
}

func (p *Parser) nextUint16() (uint16, error) {
	var res uint16
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextUint32() (uint32, error) {
	var res uint32
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextFloat32() (float32, error) {
	var res float32
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (p *Parser) nextString(l int) (string, error) {
	res := make([]byte, l)
	err := binary.Read(p.reader, p.byteOrder, &res)
	if err != nil {
		return "", err
	}
	if p.modifyStr {
		return modifyString(res), nil
	}
	return string(res), nil
}

func (p *Parser) nextFloat32Slice(l int) ([]float32, error) {
	res := make([]float32, l)
	for i := range res {
		v, err := p.nextFloat32()
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func (p *Parser) mustNextUint16() uint16 {
	res, err := p.nextUint16()
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextUint32() uint32 {
	res, err := p.nextUint32()
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextFloat32() float32 {
	res, err := p.nextFloat32()
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextFloat32Slice(l int) []float32 {
	res, err := p.nextFloat32Slice(l)
	if err != nil {
		panic(err)
	}
	return res
}

func (p *Parser) mustNextString(l int) string {
	res, err := p.nextString(l)
	if err != nil {
		panic(err)
	}
	return res
}

// modifyString 将bytes转为string，并移除末尾的空字符
func modifyString(bs []byte) string {
	i := len(bs) - 1
	for i >= 0 {
		if bs[i] != 0 {
			break
		}
		i--
	}
	return string(bs[:i+1])
}
