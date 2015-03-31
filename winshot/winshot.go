package winshot

/*
#cgo LDFLAGS: -lgdi32

#include <stdio.h>
#include <windows.h>
#include <gdiplus.h>
#include <time.h>

BITMAP bmp;
LONG Height, Width, Stride;
HDC scrdc, memdc;
HBITMAP membit;
BITMAPFILEHEADER   bmfHeader;    
BITMAPINFOHEADER   bi;

char* data;

void Raw();

void
Init() {
	scrdc = GetDC(0);
	Height = GetSystemMetrics(SM_CYSCREEN);
	Width = GetSystemMetrics(SM_CXSCREEN);
	memdc = CreateCompatibleDC(scrdc);
	membit = CreateCompatibleBitmap(scrdc, Width, Height);
	HBITMAP hOldBitmap =(HBITMAP) SelectObject(memdc, membit);
	Raw();
	Stride = bmp.bmWidthBytes;

    bi.biSize = sizeof(BITMAPINFOHEADER);    
    bi.biWidth = Width;    
    bi.biHeight = Height;  
    bi.biPlanes = 1;    
    bi.biBitCount = 32;    
    bi.biCompression = BI_RGB;    
    bi.biSizeImage = 0;  
    bi.biXPelsPerMeter = 0;    
    bi.biYPelsPerMeter = 0;    
    bi.biClrUsed = 0;    
    bi.biClrImportant = 0;

	DWORD dwBmpSize = Stride * Height;
	HANDLE hDIB = GlobalAlloc(GHND,dwBmpSize); 
    data = (char *)GlobalLock(hDIB);
}

void
Raw() {	
	BitBlt(memdc, 0, 0, Width, Height, scrdc, 0, 0, SRCCOPY);		
	int ret = GetObject(membit, sizeof(BITMAP), &bmp);
    GetDIBits(scrdc, membit, 0, (UINT)Height, data, (BITMAPINFO *)&bi, DIB_RGB_COLORS);
	printf("f: %d. %d\n", data, ret);
}

void
Clean() {
}

*/
import "C"
import "unsafe"
import "fmt"

func init() {
	C.Init()
}

func Raw() (data []byte, width, height, stride int) {
	C.Raw()
	width = int(C.Width)
	height = int(C.Height)
	stride = int(C.Stride)
	fmt.Println(C.data)
	data = (*[1 << 30]byte)(unsafe.Pointer(C.data))[0:height*stride]
	// data = nil
	return
}

func Clean() {
	C.Clean()
}
