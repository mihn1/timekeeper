//go:build windows

package browsers

import (
	"syscall"
	"unsafe"

	ole "github.com/go-ole/go-ole"
)

// ---------------------------------------------------------------------------
// Minimal UI Automation COM bindings (go-ole only)
// ---------------------------------------------------------------------------

// COM GUIDs
var (
	clsidCUIAutomation = &ole.GUID{Data1: 0xff48dba4, Data2: 0x60ef, Data3: 0x4201, Data4: [8]byte{0xaa, 0x87, 0x54, 0x10, 0x3e, 0xef, 0x59, 0x4e}}
	iidIUIAutomation   = &ole.GUID{Data1: 0x30cbe57d, Data2: 0xd9d0, Data3: 0x452a, Data4: [8]byte{0xab, 0x13, 0x7a, 0xc5, 0xac, 0x48, 0x25, 0xee}}
)

// UIA property / pattern / tree-scope constants
const (
	uiaControlTypePropertyID = 30003
	uiaNamePropertyID        = 30005
	uiaValuePatternID        = 10002
	uiaTreeScopeDescendants  = 0x4
	uiaControlTypeEdit       = 50004
)

// ---------------------------------------------------------------------------
// IUIAutomation vtable
// ---------------------------------------------------------------------------

type iuiAutomationVtbl struct {
	ole.IUnknownVtbl
	CompareElements             uintptr
	CompareRuntimeIds           uintptr
	GetRootElement              uintptr
	ElementFromHandle           uintptr
	ElementFromPoint            uintptr
	GetFocusedElement           uintptr
	GetRootElementBuildCache    uintptr
	ElementFromHandleBuildCache uintptr
	ElementFromPointBuildCache  uintptr
	GetFocusedElementBuildCache uintptr
	CreateTreeWalker            uintptr
	Get_ControlViewWalker       uintptr
	Get_ContentViewWalker       uintptr
	Get_RawViewWalker           uintptr
	Get_RawViewCondition        uintptr
	Get_ControlViewCondition    uintptr
	Get_ContentViewCondition    uintptr
	CreateCacheRequest          uintptr
	CreateTrueCondition         uintptr
	CreateFalseCondition        uintptr
	CreatePropertyCondition     uintptr
	// remaining entries omitted — we don't use them
}

type iuiAutomation struct{ ole.IUnknown }

func (v *iuiAutomation) vtbl() *iuiAutomationVtbl {
	return (*iuiAutomationVtbl)(unsafe.Pointer(v.RawVTable))
}

func newUIAutomation() (*iuiAutomation, error) {
	unk, err := ole.CreateInstance(clsidCUIAutomation, iidIUIAutomation)
	if err != nil {
		return nil, err
	}
	return (*iuiAutomation)(unsafe.Pointer(unk)), nil
}

func (v *iuiAutomation) elementFromHandle(hwnd uintptr) (*iuiAutomationElement, error) {
	var elem *iuiAutomationElement
	hr, _, _ := syscall.SyscallN(
		v.vtbl().ElementFromHandle,
		uintptr(unsafe.Pointer(v)),
		hwnd,
		uintptr(unsafe.Pointer(&elem)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return elem, nil
}

func (v *iuiAutomation) createPropertyCondition(propID uintptr, value ole.VARIANT) (*iuiAutomationCondition, error) {
	var cond *iuiAutomationCondition
	hr, _, _ := syscall.SyscallN(
		v.vtbl().CreatePropertyCondition,
		uintptr(unsafe.Pointer(v)),
		propID,
		uintptr(unsafe.Pointer(&value)),
		uintptr(unsafe.Pointer(&cond)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return cond, nil
}

// ---------------------------------------------------------------------------
// IUIAutomationCondition (opaque handle)
// ---------------------------------------------------------------------------

type iuiAutomationCondition struct{ ole.IUnknown }

// ---------------------------------------------------------------------------
// IUIAutomationElement vtable
// ---------------------------------------------------------------------------

type iuiAutomationElementVtbl struct {
	ole.IUnknownVtbl
	SetFocus                uintptr
	GetRuntimeId            uintptr
	FindFirst               uintptr
	FindAll                 uintptr
	FindFirstBuildCache     uintptr
	FindAllBuildCache       uintptr
	BuildUpdatedCache       uintptr
	GetCurrentPropertyValue uintptr
	_pad0                   uintptr // GetCurrentPropertyValueEx
	_pad1                   uintptr // GetCachedPropertyValue
	_pad2                   uintptr // GetCachedPropertyValueEx
	_pad3                   uintptr // GetCurrentPatternAs
	_pad4                   uintptr // GetCachedPatternAs
	GetCurrentPattern       uintptr
	// remaining entries omitted
}

type iuiAutomationElement struct{ ole.IUnknown }

func (v *iuiAutomationElement) vtbl() *iuiAutomationElementVtbl {
	return (*iuiAutomationElementVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *iuiAutomationElement) findFirst(scope uintptr, cond *iuiAutomationCondition) (*iuiAutomationElement, error) {
	var found *iuiAutomationElement
	hr, _, _ := syscall.SyscallN(
		v.vtbl().FindFirst,
		uintptr(unsafe.Pointer(v)),
		scope,
		uintptr(unsafe.Pointer(cond)),
		uintptr(unsafe.Pointer(&found)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return found, nil
}

func (v *iuiAutomationElement) findAll(scope uintptr, cond *iuiAutomationCondition) (*iuiAutomationElementArray, error) {
	var arr *iuiAutomationElementArray
	hr, _, _ := syscall.SyscallN(
		v.vtbl().FindAll,
		uintptr(unsafe.Pointer(v)),
		scope,
		uintptr(unsafe.Pointer(cond)),
		uintptr(unsafe.Pointer(&arr)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return arr, nil
}

func (v *iuiAutomationElement) getCurrentPattern(patternID uintptr) (*ole.IUnknown, error) {
	var pat *ole.IUnknown
	hr, _, _ := syscall.SyscallN(
		v.vtbl().GetCurrentPattern,
		uintptr(unsafe.Pointer(v)),
		patternID,
		uintptr(unsafe.Pointer(&pat)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return pat, nil
}

// ---------------------------------------------------------------------------
// IUIAutomationElementArray vtable
// ---------------------------------------------------------------------------

type iuiAutomationElementArrayVtbl struct {
	ole.IUnknownVtbl
	Get_Length uintptr
	GetElement uintptr
}

type iuiAutomationElementArray struct{ ole.IUnknown }

func (v *iuiAutomationElementArray) vtbl() *iuiAutomationElementArrayVtbl {
	return (*iuiAutomationElementArrayVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *iuiAutomationElementArray) length() (int32, error) {
	var n int32
	hr, _, _ := syscall.SyscallN(
		v.vtbl().Get_Length,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&n)),
	)
	if hr != 0 {
		return 0, ole.NewError(hr)
	}
	return n, nil
}

func (v *iuiAutomationElementArray) element(index int32) (*iuiAutomationElement, error) {
	var el *iuiAutomationElement
	hr, _, _ := syscall.SyscallN(
		v.vtbl().GetElement,
		uintptr(unsafe.Pointer(v)),
		uintptr(index),
		uintptr(unsafe.Pointer(&el)),
	)
	if hr != 0 {
		return nil, ole.NewError(hr)
	}
	return el, nil
}

// ---------------------------------------------------------------------------
// IUIAutomationValuePattern vtable
// ---------------------------------------------------------------------------

type iuiAutomationValuePatternVtbl struct {
	ole.IUnknownVtbl
	SetValue         uintptr
	Get_CurrentValue uintptr
	// remaining entries omitted
}

type iuiAutomationValuePattern struct{ ole.IUnknown }

func (v *iuiAutomationValuePattern) vtbl() *iuiAutomationValuePatternVtbl {
	return (*iuiAutomationValuePatternVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *iuiAutomationValuePattern) currentValue() (string, error) {
	var bstr *uint16
	hr, _, _ := syscall.SyscallN(
		v.vtbl().Get_CurrentValue,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&bstr)),
	)
	if hr != 0 {
		return "", ole.NewError(hr)
	}
	return ole.BstrToString(bstr), nil
}

// ---------------------------------------------------------------------------
// VARIANT helpers
// ---------------------------------------------------------------------------

func variantString(s string) ole.VARIANT {
	return ole.NewVariant(ole.VT_BSTR, int64(uintptr(unsafe.Pointer(ole.SysAllocStringLen(s)))))
}

func variantInt(i int64) ole.VARIANT {
	return ole.NewVariant(ole.VT_INT, i)
}

// ---------------------------------------------------------------------------
// Address bar names for different browsers (English locale)
// ---------------------------------------------------------------------------

var addressBarNames = []string{
	"Address and search bar",              // Chrome, Edge
	"Search or enter web address",         // Firefox
	"Search with Google or enter address", // Chrome variant
}

// ---------------------------------------------------------------------------
// Public extraction method
// ---------------------------------------------------------------------------

// extractURLUsingUIA uses the Windows UI Automation COM API to find the
// browser address bar element and read its value. This is the recommended
// approach for modern Chromium-based browsers.
func (e *BrowserURLExtractor) extractURLUsingUIA(hwnd uintptr) string {
	// Ensure COM is initialized on this thread (STA for message-pump compat).
	// CoInitialize is idempotent — returns S_FALSE if already initialized.
	ole.CoInitialize(0)
	// Do NOT call CoUninitialize — we're on the observer's message-pump
	// thread and other components depend on COM remaining active.

	auto, err := newUIAutomation()
	if err != nil {
		e.logger.Debug("UIA: NewUIAutomation failed", "error", err)
		return ""
	}
	defer auto.Release()

	root, err := auto.elementFromHandle(hwnd)
	if err != nil || root == nil {
		e.logger.Debug("UIA: ElementFromHandle failed", "error", err, "hwnd", hwnd)
		return ""
	}
	defer root.Release()

	// Strategy 1: search by known address-bar name
	for _, name := range addressBarNames {
		if url := e.uiaFindValueByName(auto, root, name); url != "" {
			return url
		}
	}

	// Strategy 2: scan all Edit descendants for a URL value
	return e.uiaFindURLInEdits(auto, root)
}

func (e *BrowserURLExtractor) uiaFindValueByName(auto *iuiAutomation, root *iuiAutomationElement, name string) string {
	cond, err := auto.createPropertyCondition(uiaNamePropertyID, variantString(name))
	if err != nil {
		return ""
	}
	defer cond.Release()

	el, err := root.findFirst(uiaTreeScopeDescendants, cond)
	if err != nil || el == nil {
		return ""
	}
	defer el.Release()

	return e.uiaReadValue(el)
}

func (e *BrowserURLExtractor) uiaFindURLInEdits(auto *iuiAutomation, root *iuiAutomationElement) string {
	cond, err := auto.createPropertyCondition(uiaControlTypePropertyID, variantInt(uiaControlTypeEdit))
	if err != nil {
		return ""
	}
	defer cond.Release()

	arr, err := root.findAll(uiaTreeScopeDescendants, cond)
	if err != nil || arr == nil {
		return ""
	}
	defer arr.Release()

	n, err := arr.length()
	if err != nil {
		return ""
	}

	for i := int32(0); i < n; i++ {
		el, err := arr.element(i)
		if err != nil || el == nil {
			continue
		}
		val := e.uiaReadValue(el)
		el.Release()
		if val != "" {
			return val
		}
	}
	return ""
}

func (e *BrowserURLExtractor) uiaReadValue(el *iuiAutomationElement) string {
	patUnk, err := el.getCurrentPattern(uiaValuePatternID)
	if err != nil || patUnk == nil {
		return ""
	}
	defer patUnk.Release()

	pat := (*iuiAutomationValuePattern)(unsafe.Pointer(patUnk))
	val, err := pat.currentValue()
	if err != nil || val == "" {
		return ""
	}

	normalized := normalizeURL(val)
	if isValidURL(normalized) {
		return normalized
	}
	return ""
}
