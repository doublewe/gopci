package gopci

import (
	"fmt"

	"github.com/hertg/gopci/pkg/pci"
	"github.com/jaypipes/pcidb"
)

type device pci.Device

type PCI struct {
	devices   []*pci.Device
	pcidbData *pcidb.PCIDB

	subclassName map[string]string
	vendorName   map[string]*pcidb.Vendor
}

func NewPCI() (*PCI, error) {
	dbData, err := pcidb.New()
	if err != nil {
		return nil, err
	}

	d, err := pci.Scan()
	if err != nil {
		return nil, err
	}

	return &PCI{
		devices:   d,
		pcidbData: dbData,
	}, nil
}

func (p *PCI) GetDevicesByHex() map[string][]*device {
	m := make(map[string][]*device)
	for _, d := range p.devices {
		m[d.Class.Hex] = append(m[d.Class.Hex], (*device)(d))
	}

	return m
}

func (p *PCI) GetDeviceByClassName(className string) []*device {
	subclassName := p.GetSubclassName()
	for classHex, d := range p.GetDevicesByHex() {
		if name, ok := subclassName[classHex]; ok && name == className {
			return d
		}
	}

	return nil
}

// 获取class对应的名称
func (p *PCI) GetSubclassName() map[string]string {
	if p.subclassName != nil {
		return p.subclassName
	}

	classNameM := make(map[string]string)
	for _, devClass := range p.pcidbData.Classes {
		for _, devSubclass := range devClass.Subclasses {
			if len(devSubclass.ProgrammingInterfaces) != 0 {
				for _, progIface := range devSubclass.ProgrammingInterfaces {
					hex := fmt.Sprintf("0x%s%s%s", devClass.ID, devSubclass.ID, progIface.ID)
					classNameM[hex] = progIface.Name
				}
			} else {
				hex := fmt.Sprintf("0x%s%s00", devClass.ID, devSubclass.ID)
				classNameM[hex] = devSubclass.Name

			}
		}
	}

	p.subclassName = classNameM
	return classNameM
}

// 获取vendor对应名称
func (p *PCI) ToVendorName(vendorHex string) string {
	if p.vendorName == nil {
		vendorNameM := make(map[string]*pcidb.Vendor)
		for _, vendor := range p.pcidbData.Vendors {
			hex := fmt.Sprintf("0x%s", vendor.ID)
			vendorNameM[hex] = vendor
		}

		p.vendorName = vendorNameM
	}

	if v, ok := p.vendorName[vendorHex]; ok {
		return v.Name
	}

	return ""
}

// 获取subvendor对应名称
func (p *PCI) ToSubVendorName(vendorHex string, productHex string, subvendorHex string) string {
	if p.vendorName == nil {
		vendorNameM := make(map[string]*pcidb.Vendor)
		for _, vendor := range p.pcidbData.Vendors {
			hex := fmt.Sprintf("0x%s", vendor.ID)
			vendorNameM[hex] = vendor
		}

		p.vendorName = vendorNameM
	}

	if vendor, ok := p.vendorName[vendorHex]; ok {
		for _, product := range vendor.Products {
			if fmt.Sprintf("0x%s", product.ID) == productHex {
				for _, subsystem := range product.Subsystems {
					if fmt.Sprintf("0x%s", subsystem.ID) == subvendorHex {
						return subsystem.Name
					}
				}
			}
		}
	}

	return ""
}

// 获取product对应名称
func (p *PCI) ToProductName(vendorHex string, productHex string) string {
	if len(p.vendorName) == 0 {
		vendorNameM := make(map[string]*pcidb.Vendor)
		for _, vendor := range p.pcidbData.Vendors {
			hex := fmt.Sprintf("0x%s", vendor.ID)
			vendorNameM[hex] = vendor
		}

		p.vendorName = vendorNameM
	}

	if vendor, ok := p.vendorName[vendorHex]; ok {
		for _, product := range vendor.Products {
			if fmt.Sprintf("0x%s", product.ID) == productHex {
				return product.Name
			}
		}
	}

	return ""
}
