package main

//#cgo pkg-config: libavformat libavcodec libavutil libswresample libavdevice libavfilter libswscale
//#include <stdio.h>
//#include <stdlib.h>
//#include <inttypes.h>
//#include <stdint.h>
//#include <string.h>
//#include <libavformat/avformat.h>
//#include <libavcodec/avcodec.h>
//#include <libavutil/avutil.h>
//#include <libavdevice/avdevice.h>
//#include <libswresample/swresample.h>
//#include <libswscale/swscale.h>
//#include <libavfilter/avfilter.h>
import "C"
import (
	"fmt"
	"log"
)

// AVError represents the C style int return code.
type AVError C.int

func (i AVError) Error() string {
	return fmt.Sprintf("av error %d", int(i))
}

func run(ret C.int) error {
	if ret != 0 {
		return AVError(ret)
	}
	return nil
}

func main() {
	log.Printf("AvFormat Version:\t%v", uint(C.avformat_version()))
	log.Printf("AvFilter Version:\t%v", uint(C.avfilter_version()))
	log.Printf("AvDevice Version:\t%v", uint(C.avdevice_version()))
	log.Printf("SWScale Version:\t%v", uint(C.swscale_version()))
	log.Printf("AvUtil Version:\t%v", uint(C.avutil_version()))
	log.Printf("AvCodec Version:\t%v", uint(C.avcodec_version()))
	log.Printf("Resample Version:\t%v", uint(C.swresample_version()))

	filename := "./test.mp4"

	C.av_register_all()

	var ctx *C.struct_AVFormatContext
	var format *C.struct_AVInputFormat

	if err := run(C.avformat_open_input(&ctx, C.CString(filename), format, nil)); err != nil {
		log.Panicln(err)
	}

	if err := run(C.avformat_find_stream_info(ctx, nil)); err != nil {
		log.Panicln(err)
	}

	C.av_dump_format(ctx, C.int(0), C.CString(filename), C.int(0))
}
