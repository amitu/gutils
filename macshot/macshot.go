package macshot

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Foundation
#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>

NSData *data;
int width, height, stride; // stride = bytes per row

void
Raw() {
	CGImageRef image = CGDisplayCreateImage(kCGDirectMainDisplay);
	width = CGImageGetWidth(image);
	height = CGImageGetHeight(image);
	stride = CGImageGetBytesPerRow(image);

	CFDataRef dataref = CGDataProviderCopyData(CGImageGetDataProvider(image));
	data = [NSData dataWithData:(NSData *)dataref];
	CFRelease(dataref);
	CGImageRelease(image);
}

void
JPEG(float quality) {
	CGImageRef image = CGDisplayCreateImage(kCGDirectMainDisplay);
	CFMutableDataRef mutableData = CFDataCreateMutable(NULL, 0);
	CGImageDestinationRef idst = CGImageDestinationCreateWithData(
		mutableData, kUTTypeJPEG, 1, NULL
	);

	NSInteger exif             =       1;
	CGFloat compressionQuality = quality;

	NSDictionary *props = [
		[NSDictionary alloc]
		initWithObjectsAndKeys:[NSNumber numberWithFloat:compressionQuality],
		kCGImageDestinationLossyCompressionQuality,
		[NSNumber numberWithInteger:exif],
		kCGImagePropertyOrientation, nil
	];

	CGImageDestinationAddImage(idst, image, (CFDictionaryRef)props);
	CGImageDestinationFinalize(idst);
	data = [NSData dataWithData:(NSData *)mutableData];
	[props release];
	CFRelease(idst);
	CFRelease(mutableData);
	CGImageRelease(image);
}

int
Length() {
	return [data length];
}

const void*
Data() {
	return [data bytes];
}

void
Free () {
	[data release];
}
*/
import "C"

import (
	"sync"
	"unsafe"
)

var (
	mutext sync.Mutex
)

// ScreenShot takes a screenshot and returns it in jpeg format with given
// quality. This is not threadsafe.
func ScreenShot(quality float64) ([]byte, error) {
	C.JPEG(C.float(quality))
	data := (*[1 << 30]byte)(unsafe.Pointer(C.Data()))[0:C.Length()]
	newData := make([]byte, C.Length())
	copy(newData, data)
	C.Free()
	return newData, nil
}

// SafeScreenShot takes a screenshot and returns it in jpeg format with given
// quality. This version is protected by a mutex and is thread safe.
func SafeScreenShot(quality float64) ([]byte, error) {
	mutext.Lock()
	defer mutext.Unlock()

	return ScreenShot(quality)
}

func Raw() (data []byte, width, height, stride int, err error) {
	C.Raw()
	// data = (*[1 << 30]byte)(unsafe.Pointer(C.Data()))[0:C.Length()]
	// newData := make([]byte, C.Length())
	// copy(newData, data)
	// C.Free()
	// return newData, int(C.width), int(C.height), int(C.stride), nil
	data = (*[1 << 30]byte)(unsafe.Pointer(C.Data()))[0:C.Length()]
	return data, int(C.width), int(C.height), int(C.stride), nil
}

func Clean() {
	C.Free()
}

/*
type ARGB struct {
	image.RGBA
}

func (p *ARGB) At(x, y int) color.Color {
	r, g, b, a := p.RGBA.At(x, y).RGBA()
	return color.RGBA{uint8(b), uint8(g), uint8(r), uint8(a)}
}

func getSS2() image.Image {
	data, w, h, s, err := macshot.Raw()
	if err != nil {
		panic(err)
	}
	return &ARGB{
		image.RGBA{
			Pix:    []uint8(data),
			Stride: s,
			Rect:   image.Rect(0, 0, w, h),
		},
	}
}

func t1() {
	f, err := os.Create("raw.jpg")
	if err != nil {
		panic(err)
	}
	jpeg.Encode(f, getSS2(), nil)
	f.Close()
}
*/
