package updater

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"news-aggregator/entity/source"
	"news-aggregator/storage/mock_aggregator"
)

var _ = Describe("News Updater", func() {
	var (
		Storage *client.MockStorage
		service Service
		logHook *test.Hook
	)

	BeforeEach(func() {
		ctrl := gomock.NewController(GinkgoT())
		Storage = client.NewMockStorage(ctrl)
		service = Service{Storage: Storage}
		logHook = test.NewGlobal()

	})

	AfterEach(func() {
		logHook.Reset()
	})
	Context("Negative cases for UpdateNews method", func() {
		It("UpdateNews should returns and log error when GetSources return error", func() {
			Storage.EXPECT().GetSources().Return(nil, errors.New("storage errors"))
			service.UpdateNews()

			Expect(logHook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
			Expect(logHook.LastEntry().Message).To(ContainSubstring("Failed to retrieve sources"))
		})

		It("should log error when updateSourceNews fails", func() {
			testSources := []source.Source{
				{
					Name:       "Test Source",
					Link:       "http://example.com",
					SourceType: source.STORAGE,
				},
			}

			Storage.EXPECT().GetSources().Return(testSources, nil)

			service.UpdateNews()

			entries := logHook.AllEntries()

			found := false
			for _, entry := range entries {
				if entry.Level == logrus.ErrorLevel && entry.Message == "Failed to update news for source: Test Source" {
					found = true
					break
				}
			}

			Expect(found).To(BeTrue(), "Expected log entry with error message not found")
		})
	})
	Context("Negative cases for updateSourceNews method", func() {
		It("updateSourceNews should returns error when GetRssFeedLink return err", func() {
			testSource := source.Source{
				Name:       "Test Source",
				Link:       "",
				SourceType: source.STORAGE,
			}
			err := updateSourceNews(testSource, nil)
			Expect(err).To(HaveOccurred())
		})

		It("updateSourceNews should returns error when SaveNews return err", func() {

			testSource := source.Source{
				Name:       "Test Source",
				Link:       "https://www.cbsnews.com/world/",
				SourceType: source.STORAGE,
			}

			Storage.EXPECT().SaveNews(gomock.Any(), gomock.Any()).Return(source.Source{Name: "pravda"}, errors.New("storage errors"))

			err := updateSourceNews(testSource, Storage)
			Expect(err).To(HaveOccurred())
		})
	})
})
