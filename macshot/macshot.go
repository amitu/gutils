package macshot

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Foundation
#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>

NSData *data;

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
import "C" // Do not merge these imports in one statement
import "unsafe"
import "sync"

var (
	screenshotMutext sync.Mutex
)

func ScreenShot(quality float64) ([]byte, error) {
	C.JPEG(C.float(quality))
	data := (*[1<<30]byte)(unsafe.Pointer(C.Data()))[0:C.Length()]
	newData := make([]byte, C.Length())
	copy(newData, data)
	C.Free()
	return newData, nil
}

func SafeScreenShot(quality float64) ([]byte, error) {
	screenshotMutext.Lock()
	defer screenshotMutext.Unlock()

	return ScreenShot(quality)
}
