package kcp_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SarunasBucius/kafka-cass-practise/kcp"
	"github.com/SarunasBucius/kafka-cass-practise/mocks"
)

func TestProduceVisit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockProdEvent := mocks.NewMockProducer(mockCtrl)
	k := kcp.New(mockProdEvent, nil)

	param := "ip"
	mockProdEvent.EXPECT().ProduceEvent(gomock.Any()).
		Return(nil).
		Times(1).
		Do(func(arg kcp.Event) {
			if param != arg.IP {
				t.Fail()
			}
		})
	k.ProduceVisit(param)

}
