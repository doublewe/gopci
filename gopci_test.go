package gopci_test

import (
	"fmt"
	"testing"

	"github.com/doublewe/gopci"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := gopci.NewPCI()
	assert.Equal(t, err, nil)
}

func TestGetDevicesByHex(t *testing.T) {
	pci, err := gopci.NewPCI()
	assert.Equal(t, err, nil)

	res := pci.GetDevicesByHex()
	for k, v := range res {
		fmt.Println(k)

		for _, v2 := range v {
			fmt.Println(v2)
		}
	}
}

func TestGetDeviceBySlot(t *testing.T) {
	pci, err := gopci.NewPCI()
	assert.Equal(t, err, nil)

	device := pci.GetDeviceBySlot("0000:00:02.0")
	fmt.Println(device)
}
