package vast

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/llgoer/go-xml/xmltree"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/stretchr/testify/assert"
)

func TestQuickStartComplex(t *testing.T) {
	skip := Duration(5 * time.Second)
	v := VAST{
		Version: "4.2",
		Ads: []Ad{
			{
				ID:            "123",
				Type:          "front",
				AdType:        "video",
				ConditionalAd: false,
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  PlainString{CDATA: "adTitle"},
					Impressions: []Impression{
						{ID: "11111", URI: "http://impressionv1.track.com"},
						{ID: "11112", URI: "http://impressionv2.track.com"},
					},
					Category: &[]Category{
						{Authority: "https://www.iabtechlab.com/categoryauthority", Category: "American Cuisine"},
						{Authority: "https://www.iabtechlab.com/categoryauthority", Category: "Guitar"},
					},
					Description: &CDATAString{"123"},
					ViewableImpression: &ViewableImpression{
						ID: "1234",
						Viewable: []CDATAString{
							{CDATA: "http://viewable1.track.com"},
							{CDATA: "http://viewable2.track.com"},
						},
						NotViewable: []CDATAString{
							{CDATA: "http://notviewable1.track.com"},
							{CDATA: "http://notviewable2.track.com"},
						},
						ViewUndetermined: []CDATAString{
							{CDATA: "http://viewundetermined1.track.com"},
							{CDATA: "http://viewundetermined2.track.com"},
						},
					},
					Creatives: []Creative{
						{
							ID:       "987",
							Sequence: 0,
							AdID:     "12",
							UniversalAdID: &[]UniversalAdID{
								{
									IDRegistry: "Ad-ID",
									ID:         "8465",
								},
								{
									IDRegistry: "FOO-ID",
									ID:         "6666465",
								},
							},
							Linear: &Linear{
								SkipOffset: &Offset{
									Duration: &skip,
								},
								Duration: Duration(15 * time.Second),
								TrackingEvents: &TrackingEvents{Tracking: []Tracking{
									{Event: EventTypeStart, URI: "http://track.xxx.com/q/start?xx"},
									{Event: EventTypeFirstQuartile, URI: "http://track.xxx.com/q/firstQuartile?xx"},
									{Event: EventTypeMidpoint, URI: "http://track.xxx.com/q/midpoint?xx"},
									{Event: EventTypeThirdQuartile, URI: "http://track.xxx.com/q/thirdQuartile?xx"},
									{Event: EventTypeComplete, URI: "http://track.xxx.com/q/complete?xx"},
								},
								},
								MediaFiles: &MediaFiles{
									MediaFile: []MediaFile{
										{
											Delivery: "progressive",
											Type:     "video/mp4",
											Width:    1024,
											Height:   576,
											URI:      "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4",
										},
									},
								},
							},
						},
					},
					Extensions: &[]Extension{
						{
							Type: "ClassName",
							Data: "AdsVideoView",
						},
						{
							Type: "ExtURL",
							Data: "http://xxxxxxxx",
						},
					},
				},
			},
		},
	}

	out, _ := xml.MarshalIndent(v, " ", "  ")
	fmt.Println(string(out))
}

func TestQuickStart(t *testing.T) {
	d := Duration(5 * time.Second)
	v := VAST{
		Mute:    true,
		Version: "3.0",
		Ads: []Ad{
			{
				ID:   "123",
				Type: "front",
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  PlainString{CDATA: "adTitle"},
					Impressions: []Impression{
						{ID: "11111", URI: "http://impressionv1.track.com"},
						{ID: "11112", URI: "http://impressionv2.track.com"},
					},
					Creatives: []Creative{
						{
							ID:       "987",
							Sequence: 0,
							Linear: &Linear{
								SkipOffset: &Offset{
									Duration: &d,
								},
								Duration: Duration(15 * time.Second),
								TrackingEvents: &TrackingEvents{Tracking: []Tracking{
									{Event: EventTypeStart, URI: "http://track.xxx.com/q/start?xx"},
									{Event: EventTypeFirstQuartile, URI: "http://track.xxx.com/q/firstQuartile?xx"},
									{Event: EventTypeMidpoint, URI: "http://track.xxx.com/q/midpoint?xx"},
									{Event: EventTypeThirdQuartile, URI: "http://track.xxx.com/q/thirdQuartile?xx"},
									{Event: EventTypeComplete, URI: "http://track.xxx.com/q/complete?xx"},
								},
								},
								MediaFiles: &MediaFiles{
									MediaFile: []MediaFile{
										{
											Delivery: "progressive",
											Type:     "video/mp4",
											Width:    1024,
											Height:   576,
											URI:      "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4",
											Label:    "123",
										},
									},
								},
							},
						},
					},
					Extensions: &[]Extension{
						{
							Type: "ClassName",
							Data: "AdsVideoView",
						},
						{
							Type: "ExtURL",
							Data: "http://xxxxxxxx",
						},
					},
				},
			},
		},
	}

	want := []byte(`{"Version":"3.0","Ad":[{"InLine":{"AdSystem":{"Data":"DSP"},"Extensions":[{"Type":"ClassName","Data":"AdsVideoView"},{"Type":"ExtURL","Data":"http://xxxxxxxx"}],"Impressions":[{"ID":"11111","URI":"http://impressionv1.track.com"},{"ID":"11112","URI":"http://impressionv2.track.com"}],"AdTitle":{"Data":"adTitle"},"Creatives":[{"ID":"987","Linear":{"SkipOffset":"00:00:05","Duration":"00:00:15","TrackingEvents":{"Tracking":[{"Event":"start","URI":"http://track.xxx.com/q/start?xx"},{"Event":"firstQuartile","URI":"http://track.xxx.com/q/firstQuartile?xx"},{"Event":"midpoint","URI":"http://track.xxx.com/q/midpoint?xx"},{"Event":"thirdQuartile","URI":"http://track.xxx.com/q/thirdQuartile?xx"},{"Event":"complete","URI":"http://track.xxx.com/q/complete?xx"}]},"MediaFiles":{"MediaFile":[{"Delivery":"progressive","Type":"video/mp4","Width":1024,"Height":576,"URI":"http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4","Label":"123"}]}}}]},"ID":"123","Type":"front"}],"Mute":true}`)
	got, err := json.Marshal(v)
	t.Logf("%s", got)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n got = %s,  \nwant = %s", got, want)
	}

}

func TestEmptyVast(t *testing.T) {
	v := VAST{
		Version: "3.0",
		Errors: []CDATAString{
			{CDATA: "http://xx.xx.com/e/error?e=__ERRORCODE__&co=__CONTENTPLAYHEAD__&ca=__CACHEBUSTING__&a=__ASSETURI__&t=__TIMESTAMP__&o=__OTHER__"},
		},
	}
	want := []byte(`{"Version":"3.0","Errors":[{"Data":"http://xx.xx.com/e/error?e=__ERRORCODE__\u0026co=__CONTENTPLAYHEAD__\u0026ca=__CACHEBUSTING__\u0026a=__ASSETURI__\u0026t=__TIMESTAMP__\u0026o=__OTHER__"}]}`)
	got, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Logf("%s", got)
		t.Errorf("Marshal() got = %v, want %v", got, want)
	}

	want = []byte(`<VAST version="3.0"><Error><![CDATA[http://xx.xx.com/e/error?e=__ERRORCODE__&co=__CONTENTPLAYHEAD__&ca=__CACHEBUSTING__&a=__ASSETURI__&t=__TIMESTAMP__&o=__OTHER__]]></Error></VAST>`)
	got, err = xml.Marshal(v)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Logf("%s", got)
		t.Errorf("Marshal() got = %v, want %v", got, want)
	}

}

func createVastDemo() (*VAST, error) {
	adId := "123"
	adTitle := "ad title"
	assetId := "123456"
	impressionId := "456"
	impressionURI := "http://impression.track.cn"
	seconds := Duration(15 * time.Second)
	mediaType := "video/mp4"
	mediaURI := "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4"

	v := &VAST{
		Version: "3.0",
		XMLNS:   "http://www.iab.com/VAST",
		Ads: []Ad{
			{
				ID: adId,
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  PlainString{CDATA: adTitle},
					Impressions: []Impression{
						{ID: impressionId, URI: impressionURI},
					},
					Creatives: []Creative{
						{
							ID:       assetId,
							Sequence: 0,
							Linear: &Linear{
								Duration: seconds,
								TrackingEvents: &TrackingEvents{Tracking: []Tracking{
									{
										Event:  EventTypeStart,
										Offset: nil,
										URI:    "http://track.xxx.com/q/start?xx",
										UA:     "",
									},
								},
								},
								MediaFiles: &MediaFiles{
									MediaFile: []MediaFile{
										{
											Delivery: "progressive",
											Type:     mediaType,
											Width:    1024,
											Height:   576,
											URI:      mediaURI,
											Label:    "123",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return v, nil
}

func BenchmarkVastMarshalXML(b *testing.B) {

	want := []byte(`<VAST version="3.0" xmlns="http://www.iab.com/VAST"><Ad id="123" type="front"><InLine><AdSystem><![CDATA[DSP]]></AdSystem><AdTitle><![CDATA[ad title]]></AdTitle><Impression id="456"><![CDATA[http://impression.track.cn]]></Impression><Creatives><Creative id="123456"><Linear><Duration>00:00:15</Duration><TrackingEvents><Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking></TrackingEvents><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1024" height="576"><![CDATA[http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4]]></MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v, _ := createVastDemo()
		got, err := xml.Marshal(v)
		if err != nil {
			b.Errorf("Marshal() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			b.Errorf("Marshal() got = %v, want %v", got, want)
		}
	}
}

func BenchmarkVastMarshalJson(b *testing.B) {

	want := []byte(`{"Version":"3.0","xmlns":"http://www.iab.com/VAST","Ad":[{"ID":"123","Type":"front","InLine":{"AdSystem":{"Data":"DSP"},"AdTitle":{"Data":"ad title"},"Impressions":[{"ID":"456","URI":"http://impression.track.cn"}],"Creatives":[{"ID":"123456","Linear":{"Duration":"00:00:15"{"Tracking":[{"Event":"start","URI":"http://track.xxx.com/q/start?xx"}]},"MediaFiles":[{"Delivery":"progressive","Type":"video/mp4","Width":1024,"Height":576,"URI":"http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4"}]}}]}}]}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v, _ := createVastDemo()
		got, err := ffjson.Marshal(v)
		if err != nil {
			b.Errorf("Marshal() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			b.Errorf("Marshal() got = %v, want %v", got, want)
		}
	}
}

func TestCreateVastJson(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{name: "testCase1", want: []byte(`{"Version":"3.0","xmlns":"http://www.iab.com/VAST","Ad":[{"InLine":{"AdSystem":{"Data":"DSP"},"Impressions":[{"ID":"456","URI":"http://impression.track.cn"}],"AdTitle":{"Data":"ad title"},"Creatives":[{"ID":"123456","Linear":{"Duration":"00:00:15","TrackingEvents":{"Tracking":[{"Event":"start","URI":"http://track.xxx.com/q/start?xx"}]},"MediaFiles":{"MediaFile":[{"Delivery":"progressive","Type":"video/mp4","Width":1024,"Height":576,"URI":"http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4","Label":"123"}]}}}]},"ID":"123"}]}`),
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := createVastDemo()
			got, err := ffjson.Marshal(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot  %s \nwant %s", got, tt.want)
			}
		})
	}
}

func TestCreateVastXML(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{name: "testCase1", want: []byte(`<VAST version="3.0" xmlns="http://www.iab.com/VAST"><Ad id="123"><InLine><AdSystem>DSP</AdSystem><Impression id="456"><![CDATA[http://impression.track.cn]]></Impression><AdTitle>ad title</AdTitle><Creatives><Creative id="123456"><Linear><Duration>00:00:15</Duration><TrackingEvents><Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking></TrackingEvents><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1024" height="576" label="123"><![CDATA[http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4]]></MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`),
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := createVastDemo()
			got, err := xml.Marshal(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot  %s \nwant %s", got, tt.want)
			}
		})
	}
}
func loadFixture(path string) (*VAST, []byte, string, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, nil, "", err
	}
	defer xmlFile.Close()
	b, _ := io.ReadAll(xmlFile)

	var v VAST
	err = xml.Unmarshal(b, &v)

	res, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, "", err

	}

	return &v, b, string(res), err
}

func TestCreativeExtensions(t *testing.T) {
	v, _, _, err := loadFixture("testdata/creative_extensions.xml")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "3.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "abc123", ad.ID)
		if assert.NotNil(t, ad.InLine) {
			if assert.Len(t, ad.InLine.Creatives, 1) {
				exts := *ad.InLine.Creatives[0].CreativeExtensions
				if assert.Len(t, exts, 4) {
					var ext Extension
					// asserting first extension
					ext = exts[0]
					assert.Equal(t, "geo", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n              <Country>US</Country>\n              <Bandwidth>3</Bandwidth>\n              <BandwidthKbps>1680</BandwidthKbps>\n            ", string(ext.Data))
					// asserting second extension
					ext = exts[1]
					assert.Equal(t, "activeview", ext.Type)
					if assert.Len(t, ext.CustomTracking, 2) {
						// first tracker
						assert.Equal(t, "viewable_impression", ext.CustomTracking[0].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=viewable_impression&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]&ad_mt=[AD_MT]", ext.CustomTracking[0].URI)
						// second tracker
						assert.Equal(t, "abandon", ext.CustomTracking[1].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=video_abandon&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]", ext.CustomTracking[1].URI)
					}
					assert.Empty(t, string(ext.Data))
					// asserting third extension
					ext = exts[2]
					assert.Equal(t, "DFP", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n              <SkippableAdType>Generic</SkippableAdType>\n            ", string(ext.Data))
					// asserting fourth extension
					ext = exts[3]
					assert.Equal(t, "metrics", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n              <FeEventId>MubmWKCWLs_tiQPYiYrwBw</FeEventId>\n              <AdEventId>CIGpsPCTkdMCFdN-Ygod-xkCKQ</AdEventId>\n            ", string(ext.Data))
				}
			}
		}
	}
}

func TestInlineExtensions(t *testing.T) {
	v, _, _, err := loadFixture("testdata/inline_extensions.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "3.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "708365173", ad.ID)
		if assert.NotNil(t, ad.InLine) {
			if assert.NotNil(t, ad.InLine.Extensions) {
				exts := *ad.InLine.Extensions
				if assert.Len(t, exts, 4) {
					var ext Extension
					// asserting first extension
					ext = exts[0]
					assert.Equal(t, "geo", ext.Type)
					assert.Empty(t, ext.CustomTracking)
					assert.Equal(t, "\n          <Country>US</Country>\n          <Bandwidth>3</Bandwidth>\n          <BandwidthKbps>1680</BandwidthKbps>\n        ", string(ext.Data))
					// asserting second extension
					ext = exts[1]
					assert.Equal(t, "activeview", ext.Type)
					if assert.Len(t, ext.CustomTracking, 2) {
						// first tracker
						assert.Equal(t, "viewable_impression", ext.CustomTracking[0].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=viewable_impression&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]&ad_mt=[AD_MT]", ext.CustomTracking[0].URI)
						// second tracker
						assert.Equal(t, "abandon", ext.CustomTracking[1].Event)
						assert.Equal(t, "https://pubads.g.doubleclick.net/pagead/conversion/?ai=test&label=video_abandon&acvw=[VIEWABILITY]&gv=[GOOGLE_VIEWABILITY]", ext.CustomTracking[1].URI)
					}
					assert.Empty(t, string(ext.Data))
					// asserting third extension
					ext = exts[2]
					assert.Equal(t, "DFP", ext.Type)
					assert.Equal(t, "\n          <SkippableAdType>Generic</SkippableAdType>\n        ", string(ext.Data))
					assert.Empty(t, ext.CustomTracking)
					// asserting fourth extension
					ext = exts[3]
					assert.Equal(t, "metrics", ext.Type)
					assert.Equal(t, "\n          <FeEventId>MubmWKCWLs_tiQPYiYrwBw</FeEventId>\n          <AdEventId>CIGpsPCTkdMCFdN-Ygod-xkCKQ</AdEventId>\n        ", string(ext.Data))
					assert.Empty(t, ext.CustomTracking)
				}
			}
		}
	}
}

func TestInlineLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_inline_linear.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "601364", ad.ID)
		assert.Nil(t, ad.Wrapper)
		assert.Equal(t, 0, ad.Sequence)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Acudeo Compatible", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "VAST 2.0 Instream Test 1", inline.AdTitle.CDATA)
			assert.Equal(t, "VAST 2.0 Instream Test 1", inline.Description.CDATA)
			if assert.Len(t, inline.Errors, 2) {
				assert.Equal(t, "http://myErrorURL/error", inline.Errors[0].CDATA)
				assert.Equal(t, "http://myErrorURL/error2", inline.Errors[1].CDATA)
			}
			if assert.Len(t, inline.Impressions, 2) {
				assert.Equal(t, "http://myTrackingURL/impression", inline.Impressions[0].URI)
				assert.Equal(t, "http://myTrackingURL/impression2", inline.Impressions[1].URI)
				assert.Equal(t, "foo", inline.Impressions[1].ID)
			}
			if assert.Len(t, inline.Creatives, 2) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, "601364", crea1.AdID)
				assert.Nil(t, crea1.NonLinearAds)
				assert.Nil(t, crea1.CompanionAds)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					assert.Equal(t, Duration(30*time.Second), linear.Duration)
					if assert.Len(t, linear.TrackingEvents.Tracking, 6) {
						assert.Equal(t, linear.TrackingEvents.Tracking[0].Event, "creativeView")
						assert.Equal(t, linear.TrackingEvents.Tracking[0].URI, "http://myTrackingURL/creativeView")
						assert.Equal(t, linear.TrackingEvents.Tracking[1].Event, "start")
						assert.Equal(t, linear.TrackingEvents.Tracking[1].URI, "http://myTrackingURL/start")
					}
					if assert.NotNil(t, linear.VideoClicks) {
						if assert.Len(t, linear.VideoClicks.ClickThroughs, 1) {
							assert.Equal(t, linear.VideoClicks.ClickThroughs[0].URI, "http://www.tremormedia.com")
						}
						if assert.Len(t, linear.VideoClicks.ClickTrackings, 1) {
							assert.Equal(t, linear.VideoClicks.ClickTrackings[0].URI, "http://myTrackingURL/click")
						}
						assert.Len(t, linear.VideoClicks.CustomClicks, 0)
					}
					if assert.Len(t, linear.MediaFiles.MediaFile, 1) {
						mf := linear.MediaFiles.MediaFile
						assert.Equal(t, "progressive", mf[0].Delivery)
						assert.Equal(t, "video/x-flv", mf[0].Type)
						assert.Equal(t, 500, mf[0].Bitrate)
						assert.Equal(t, 400, mf[0].Width)
						assert.Equal(t, 300, mf[0].Height)
						assert.Equal(t, true, mf[0].Scalable)
						assert.Equal(t, true, mf[0].MaintainAspectRatio)
						assert.Equal(t, "http://cdnp.tremormedia.com/video/acudeo/Carrot_400x300_500kb.flv", mf[0].URI)
					}
				}

				crea2 := inline.Creatives[1]
				assert.Equal(t, "601364-Companion", crea2.AdID)
				assert.Nil(t, crea2.NonLinearAds)
				assert.Nil(t, crea2.Linear)
				if assert.NotNil(t, crea2.CompanionAds) {
					assert.Equal(t, "all", crea2.CompanionAds.Required)
					if assert.Len(t, crea2.CompanionAds.Companions, 2) {
						comp1 := crea2.CompanionAds.Companions[0]
						assert.Equal(t, 300, comp1.Width)
						assert.Equal(t, 250, comp1.Height)
						if assert.NotNil(t, comp1.StaticResource) {
							assert.Equal(t, "image/jpeg", comp1.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/Blistex1.jpg", comp1.StaticResource.URI)
						}
						if assert.Len(t, comp1.TrackingEvents.Tracking, 1) {
							assert.Equal(t, "creativeView", comp1.TrackingEvents.Tracking[0].Event)
							assert.Equal(t, "http://myTrackingURL/firstCompanionCreativeView", comp1.TrackingEvents.Tracking[0].URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp1.CompanionClickThrough.CDATA)

						comp2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 728, comp2.Width)
						assert.Equal(t, 90, comp2.Height)
						if assert.NotNil(t, comp2.StaticResource) {
							assert.Equal(t, "image/jpeg", comp2.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/728x90_banner1.jpg", comp2.StaticResource.URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp2.CompanionClickThrough.CDATA)
					}
				}
			}
		}
	}
}

func TestInlineLinearDurationUndefined(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_inline_linear-duration_undefined.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			if assert.Len(t, inline.Creatives, 1) {
				crea1 := inline.Creatives[0]
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					assert.Equal(t, Duration(0), linear.Duration)
				}
			}
		}
	}
}

func TestInlineNonLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_inline_nonlinear.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "602678", ad.ID)
		assert.Nil(t, ad.Wrapper)
		assert.Equal(t, 0, ad.Sequence)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Acudeo Compatible", inline.AdSystem.Name)
			assert.Equal(t, "NonLinear Test Campaign 1", inline.AdTitle.CDATA)
			assert.Equal(t, "NonLinear Test Campaign 1", inline.Description.CDATA)
			assert.Equal(t, "http://mySurveyURL/survey", inline.Survey.URI)
			if assert.Len(t, inline.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/error", inline.Errors[0].CDATA)
			}
			if assert.Len(t, inline.Impressions, 1) {
				assert.Equal(t, "http://myTrackingURL/impression", inline.Impressions[0].URI)
			}
			if assert.Len(t, inline.Creatives, 2) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, "602678-NonLinear", crea1.AdID)
				assert.Nil(t, crea1.Linear)
				assert.Nil(t, crea1.CompanionAds)
				if assert.NotNil(t, crea1.NonLinearAds) {
					nonlin := crea1.NonLinearAds
					if assert.Len(t, nonlin.TrackingEvents.Tracking, 5) {
						assert.Equal(t, nonlin.TrackingEvents.Tracking[0].Event, "creativeView")
						assert.Equal(t, nonlin.TrackingEvents.Tracking[0].URI, "http://myTrackingURL/nonlinear/creativeView")
						assert.Equal(t, nonlin.TrackingEvents.Tracking[1].Event, "expand")
						assert.Equal(t, nonlin.TrackingEvents.Tracking[1].URI, "http://myTrackingURL/nonlinear/expand")
					}
					if assert.Len(t, nonlin.NonLinears, 2) {
						assert.Equal(t, "image/jpeg", nonlin.NonLinears[0].StaticResource.CreativeType)
						assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/50x300_static.jpg", strings.TrimSpace(nonlin.NonLinears[0].StaticResource.URI))
						assert.Equal(t, "image/jpeg", nonlin.NonLinears[1].StaticResource.CreativeType)
						assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/50x450_static.jpg", strings.TrimSpace(nonlin.NonLinears[1].StaticResource.URI))
						assert.Equal(t, "http://www.tremormedia.com", strings.TrimSpace(nonlin.NonLinears[1].NonLinearClickThrough.CDATA))
					}
				}

				crea2 := inline.Creatives[1]
				assert.Equal(t, "602678-Companion", crea2.AdID)
				assert.Nil(t, crea2.NonLinearAds)
				assert.Nil(t, crea2.Linear)
				if assert.NotNil(t, crea2.CompanionAds) {
					if assert.Len(t, crea2.CompanionAds.Companions, 2) {
						comp1 := crea2.CompanionAds.Companions[0]
						assert.Equal(t, 300, comp1.Width)
						assert.Equal(t, 250, comp1.Height)
						if assert.NotNil(t, comp1.StaticResource) {
							assert.Equal(t, "application/x-shockwave-flash", comp1.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/300x250_companion_1.swf", comp1.StaticResource.URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp1.CompanionClickThrough.CDATA)

						comp2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 728, comp2.Width)
						assert.Equal(t, 90, comp2.Height)
						if assert.NotNil(t, comp2.StaticResource) {
							assert.Equal(t, "image/jpeg", comp2.StaticResource.CreativeType)
							assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/728x90_banner1.jpg", comp2.StaticResource.URI)
						}
						if assert.Len(t, comp2.TrackingEvents.Tracking, 1) {
							assert.Equal(t, "creativeView", comp2.TrackingEvents.Tracking[0].Event)
							assert.Equal(t, "http://myTrackingURL/secondCompanion", comp2.TrackingEvents.Tracking[0].URI)
						}
						assert.Equal(t, "http://www.tremormedia.com", comp2.CompanionClickThrough.CDATA)
					}
				}
			}
		}
	}
}

func TestWrapperLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_wrapper_linear_1.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "602833", ad.ID)
		assert.Equal(t, 0, ad.Sequence)
		assert.Nil(t, ad.InLine)
		if assert.NotNil(t, ad.Wrapper) {
			wrapper := ad.Wrapper
			assert.Equal(t, true, *wrapper.FallbackOnNoAd)
			assert.Equal(t, true, *wrapper.AllowMultipleAds)
			assert.Nil(t, wrapper.FollowAdditionalWrappers)
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_linear.xml", wrapper.VASTAdTagURI.CDATA)
			assert.Equal(t, "Acudeo Compatible", wrapper.AdSystem.Name)
			if assert.Len(t, wrapper.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/wrapper/error", wrapper.Errors[0].CDATA)
			}
			if assert.Len(t, wrapper.Impressions, 1) {
				assert.Equal(t, "http://myTrackingURL/wrapper/impression", wrapper.Impressions[0].URI)
			}

			if assert.Len(t, wrapper.Creatives, 3) {
				crea1 := wrapper.Creatives[0]
				assert.Equal(t, "602833", crea1.AdID)
				assert.Nil(t, crea1.NonLinearAds)
				assert.Nil(t, crea1.CompanionAds)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					if assert.Len(t, linear.TrackingEvents.Tracking, 11) {
						assert.Equal(t, linear.TrackingEvents.Tracking[0].Event, "creativeView")
						assert.Equal(t, linear.TrackingEvents.Tracking[0].URI, "http://myTrackingURL/wrapper/creativeView")
						assert.Equal(t, linear.TrackingEvents.Tracking[1].Event, "start")
						assert.Equal(t, linear.TrackingEvents.Tracking[1].URI, "http://myTrackingURL/wrapper/start")
					}
					assert.Nil(t, linear.VideoClicks)
				}

				crea2 := wrapper.Creatives[1]
				assert.Equal(t, "", crea2.AdID)
				assert.Nil(t, crea2.CompanionAds)
				assert.Nil(t, crea2.NonLinearAds)
				if assert.NotNil(t, crea2.Linear) {
					if assert.Len(t, crea2.Linear.VideoClicks.ClickTrackings, 1) {
						assert.Equal(t, "http://myTrackingURL/wrapper/click", crea2.Linear.VideoClicks.ClickTrackings[0].URI)
					}
				}

				crea3 := wrapper.Creatives[2]
				assert.Equal(t, "602833-NonLinearTracking", crea3.AdID)
				assert.Nil(t, crea3.CompanionAds)
				assert.Nil(t, crea3.Linear)
				if assert.NotNil(t, crea3.NonLinearAds) {
					if assert.Len(t, crea3.NonLinearAds.TrackingEvents.Tracking, 1) {
						assert.Equal(t, "creativeView", crea3.NonLinearAds.TrackingEvents.Tracking[0].Event)
						assert.Equal(t, "http://myTrackingURL/wrapper/creativeView", crea3.NonLinearAds.TrackingEvents.Tracking[0].URI)
					}
				}
			}
		}
	}
}

func TestWrapperNonLinear(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_wrapper_nonlinear_1.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "602867", ad.ID)
		assert.Equal(t, 0, ad.Sequence)
		assert.Nil(t, ad.InLine)
		if assert.NotNil(t, ad.Wrapper) {
			wrapper := ad.Wrapper
			assert.Equal(t, "http://demo.tremormedia.com/proddev/vast/vast_inline_nonlinear2.xml", wrapper.VASTAdTagURI.CDATA)
			assert.Equal(t, "Acudeo Compatible", wrapper.AdSystem.Name)
			if assert.Len(t, wrapper.Errors, 1) {
				assert.Equal(t, "http://myErrorURL/wrapper/error", wrapper.Errors[0].CDATA)
			}
			if assert.Len(t, wrapper.Impressions, 1) {
				assert.Equal(t, "http://myTrackingURL/wrapper/impression", wrapper.Impressions[0].URI)
			}

			if assert.Len(t, wrapper.Creatives, 2) {
				crea1 := wrapper.Creatives[0]
				assert.Equal(t, "602867", crea1.AdID)
				assert.Nil(t, crea1.NonLinearAds)
				assert.Nil(t, crea1.CompanionAds)
				assert.NotNil(t, crea1.Linear)

				crea2 := wrapper.Creatives[1]
				assert.Equal(t, "602867-NonLinearTracking", crea2.AdID)
				assert.Nil(t, crea2.CompanionAds)
				assert.Nil(t, crea2.Linear)
				if assert.NotNil(t, crea2.NonLinearAds) {
					if assert.Len(t, crea2.NonLinearAds.TrackingEvents.Tracking, 5) {
						assert.Equal(t, "creativeView", crea2.NonLinearAds.TrackingEvents.Tracking[0].Event)
						assert.Equal(t, "http://myTrackingURL/wrapper/nonlinear/creativeView/creativeView", crea2.NonLinearAds.TrackingEvents.Tracking[0].URI)
					}
				}
			}
		}
	}
}

func TestSpotXVpaid(t *testing.T) {
	v, _, _, err := loadFixture("testdata/spotx_vpaid.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "1130507-1818483", ad.ID)
		assert.Nil(t, ad.Wrapper)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "SpotXchange", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "IntegralAds_VAST_2_0_Ad_Wrapper", inline.AdTitle.CDATA)
			assert.Equal(t, "", inline.Description.CDATA)

			if assert.Len(t, inline.Creatives, 2) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, 1, crea1.Sequence)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear
					adParam, err := os.Open("testdata/spotx_adparameters.txt")
					if err != nil {
						assert.FailNow(t, "Cannot open adparams file")
					}
					defer adParam.Close()
					b, _ := ioutil.ReadAll(adParam)
					assert.Equal(t, linear.AdParameters.Parameters, string(b))
					if assert.Len(t, crea1.Linear.MediaFiles.MediaFile, 1) {
						media1 := crea1.Linear.MediaFiles.MediaFile
						assert.Equal(t, "progressive", media1[0].Delivery)
						assert.Equal(t, "application/javascript", media1[0].Type)
						assert.Equal(t, 300, media1[0].Width)
						assert.Equal(t, 250, media1[0].Height)
						assert.Equal(t, "https://cdn.spotxcdn.com/integration/instreamadbroker/v1/instreamadbroker/beta.js", media1[0].URI)
					}
				}
				crea2 := inline.Creatives[1]
				assert.Equal(t, 1, crea2.Sequence)
				if assert.NotNil(t, crea2.CompanionAds) {
					if assert.Len(t, crea2.CompanionAds.Companions, 3) {
						companionAds1 := crea2.CompanionAds.Companions[0]
						assert.Equal(t, 300, companionAds1.Width)
						assert.Equal(t, 250, companionAds1.Height)
						assert.Equal(t, "medium_rectangle", companionAds1.ID)
						if assert.NotNil(t, companionAds1.IFrameResource) {
							assert.Equal(t, "https://search.spotxchange.com/banner?_a=137530&amp;_p=spotx&amp;_z=1&amp;_m=eNpVz99vgjAQB%2FD%2BLTyP0rM%2FLCZ7WOKWbFF5GJr4RAotUBUwgkHZ9rdvVB8W%2B3BJv3eXT45hQhBCwAnhEkKQKA2lZCHLhJYSwExTSgSkqcpBCEIIIKBTTglaf6JQ4glleAISAxOo7LpjOwuCvu9xrip7uJaq1tdK1ThrqgB9eefWnJLibLU385hRoTQ8FylPnSNzMgUtc6EpdY73dB%2B3bVKbflwYg6xp9tYkuqmUrceoNeqUlbg9Nt0lG7HCOGkcVGdtTZ2Z5IHyneU7zHea%2F8D93A6bcPR7e2i1XizuBW4R4oJcKPDx%2B99y5TuKD%2BU23rPlrhiW849dFG%2FstnqrVsPLJYrfYTm88mi%2BZ6uheP4DQppvcA%3D%3D&amp;_l=eNplj01rwkAYhN%2FfsleLvJv9yEbooUXoxU1BDKW5yGazNUk1SZMNtVr%2Fe9VgofQyh5mHGYYqqrhiEEYykIyDABWgAi4wlAAUKKdcRIpD3zZ%2BD5MJSBpB3dQOEI6EkdmRZGVOZgSniEySO9Kar2bwNwcvVt%2B6%2BpeJrtCQbUu7dntbmHrj1m%2FOXXOKiOfYDp3xLl%2FvTPfufLs11v1n8cKeThAni8UoCDhVAparR%2FDdsGuLxrttP7XNDlBxRjORW5pxhRYjIYULc8nDyDLm8vH3%2BRxy%2BBhM7a3p%2FdhKWSgC8XcGvuNK8%2FQp2b%2BuiiKdb1AfNNfVw6eukkBXcfk8T4L0ZVnqQ3L%2FAyjAYKo%3D&amp;_t=eNotj9Fu4yAQRf0tPFdbwK5dE%2B1b1aiVnGpXlVaJKlkDTGxaYyxMtsmm%2BfeFuDzA3MO9wyBhHNFn9Z5WQCtd6vt7xrCSOa2YlLBnZUkpZZnqwYyZ8x2MRmVySZ2JRW2gNZoIorWUpeSasbq4y%2BOOulaFrrGGUtccyQ3ZO28hRO%2F66TFKayy24TRhJMZCh2%2B3ndnHi0%2BjQx9hTmlUPZquTyl%2Bd5XBDq1ydhrQ4ph4hAZka%2BwhijTRwbYeVYCxG9KzM4JXfTtPqIg4kwk8WAzo56Q0%2FjUKU4XHkI55cuHYLriFafqekF3iwE7jQAS%2FIS5m%2BSWuLGM%2FKM1AVOI8i1KQbnASBrIygq6uYHJzWACrGY8sv%2F42gbyoKlpEVAiiTDh9m9gS9NgZNyZWsjqiKprcYQz%2B6uM8X3rppVfJ2eqSfW3Wv05b%2Fvy%2Be2j4br352NqN3f55Ys2%2Fj8%2FNQ3d8WTfHhv8eXl6bn%2F8BtbOdiA%3D%3D&amp;_b=eNozYEjRB8KUpCSzJKMUQ0NLE1NjIJmaYplskmKZaplolmJplKqXUZKbw%2BAX6uODIGp83T0NI7MysiOzHCv8wj0rI3MDq3yNwnKAtEFkuKtBVDiQrvI1jqxKtwUAEJMfUw%3D%3D&amp;resource_type=iframe", companionAds1.IFrameResource.CDATA)
						}
						companionAds2 := crea2.CompanionAds.Companions[1]
						assert.Equal(t, 300, companionAds2.Width)
						assert.Equal(t, 250, companionAds2.Height)
						assert.Equal(t, "medium_rectangle", companionAds2.ID)
						if assert.NotNil(t, companionAds2.HTMLResource) {
							htmlResource, err := os.Open("testdata/spotx_html_resource.html")
							if err != nil {
								assert.FailNow(t, "Cannot open spotx html resource file")
							}
							defer htmlResource.Close()
							b, _ := ioutil.ReadAll(htmlResource)
							assert.Equal(t, companionAds2.HTMLResource.HTML, string(b))
						}
						companionAds3 := crea2.CompanionAds.Companions[2]
						assert.Equal(t, 300, companionAds3.Width)
						assert.Equal(t, 250, companionAds3.Height)
						assert.Equal(t, "medium_rectangle", companionAds3.ID)
						if assert.NotNil(t, companionAds3.StaticResource) {
							assert.Equal(t, "image/gif", companionAds3.StaticResource.CreativeType)
							assert.Equal(t, "https://search.spotxchange.com/banner?_a=137530&amp;_p=spotx&amp;_z=1&amp;_m=eNpVz99vgjAQB%2FD%2BLTyP0rM%2FLCZ7WOKWbFF5GJr4RAotUBUwgkHZ9rdvVB8W%2B3BJv3eXT45hQhBCwAnhEkKQKA2lZCHLhJYSwExTSgSkqcpBCEIIIKBTTglaf6JQ4glleAISAxOo7LpjOwuCvu9xrip7uJaq1tdK1ThrqgB9eefWnJLibLU385hRoTQ8FylPnSNzMgUtc6EpdY73dB%2B3bVKbflwYg6xp9tYkuqmUrceoNeqUlbg9Nt0lG7HCOGkcVGdtTZ2Z5IHyneU7zHea%2F8D93A6bcPR7e2i1XizuBW4R4oJcKPDx%2B99y5TuKD%2BU23rPlrhiW849dFG%2FstnqrVsPLJYrfYTm88mi%2BZ6uheP4DQppvcA%3D%3D&amp;_l=eNplj01rwkAYhN%2FfsleLvJv9yEbooUXoxU1BDKW5yGazNUk1SZMNtVr%2Fe9VgofQyh5mHGYYqqrhiEEYykIyDABWgAi4wlAAUKKdcRIpD3zZ%2BD5MJSBpB3dQOEI6EkdmRZGVOZgSniEySO9Kar2bwNwcvVt%2B6%2BpeJrtCQbUu7dntbmHrj1m%2FOXXOKiOfYDp3xLl%2FvTPfufLs11v1n8cKeThAni8UoCDhVAparR%2FDdsGuLxrttP7XNDlBxRjORW5pxhRYjIYULc8nDyDLm8vH3%2BRxy%2BBhM7a3p%2FdhKWSgC8XcGvuNK8%2FQp2b%2BuiiKdb1AfNNfVw6eukkBXcfk8T4L0ZVnqQ3L%2FAyjAYKo%3D&amp;_t=eNotj9Fu4yAQRf0tPFdbwK5dE%2B1b1aiVnGpXlVaJKlkDTGxaYyxMtsmm%2BfeFuDzA3MO9wyBhHNFn9Z5WQCtd6vt7xrCSOa2YlLBnZUkpZZnqwYyZ8x2MRmVySZ2JRW2gNZoIorWUpeSasbq4y%2BOOulaFrrGGUtccyQ3ZO28hRO%2F66TFKayy24TRhJMZCh2%2B3ndnHi0%2BjQx9hTmlUPZquTyl%2Bd5XBDq1ydhrQ4ph4hAZka%2BwhijTRwbYeVYCxG9KzM4JXfTtPqIg4kwk8WAzo56Q0%2FjUKU4XHkI55cuHYLriFafqekF3iwE7jQAS%2FIS5m%2BSWuLGM%2FKM1AVOI8i1KQbnASBrIygq6uYHJzWACrGY8sv%2F42gbyoKlpEVAiiTDh9m9gS9NgZNyZWsjqiKprcYQz%2B6uM8X3rppVfJ2eqSfW3Wv05b%2Fvy%2Be2j4br352NqN3f55Ys2%2Fj8%2FNQ3d8WTfHhv8eXl6bn%2F8BtbOdiA%3D%3D&amp;_b=eNpFxl0LgjAUgGF%2FkSc%2FEhZ0EUje6AFFibxznuXmsklbGNKPL6%2FihYd350nnZnsA6Onh29m49za9mWASpDpw8jVxC%2BapBqAt4jzhIQUBi%2FfRT0Gsj4kJ1iXEQuEP6uZhk%2Bd%2FPsVaKUxPK6Z3haOO26xZ2gvKoi6Xa30eMas0rmWEtT5%2BAZ5CLzA%3D", companionAds3.StaticResource.URI)
						}
						if assert.NotNil(t, companionAds3.CompanionClickThrough) {
							assert.Equal(t, "https://search.spotxchange.com/click?_a=137530&amp;_p=spotx&amp;_z=1&amp;_m=eNpVz99vgjAQB%2FD%2BLTyP0rM%2FLCZ7WOKWbFF5GJr4RAotUBUwgkHZ9rdvVB8W%2B3BJv3eXT45hQhBCwAnhEkKQKA2lZCHLhJYSwExTSgSkqcpBCEIIIKBTTglaf6JQ4glleAISAxOo7LpjOwuCvu9xrip7uJaq1tdK1ThrqgB9eefWnJLibLU385hRoTQ8FylPnSNzMgUtc6EpdY73dB%2B3bVKbflwYg6xp9tYkuqmUrceoNeqUlbg9Nt0lG7HCOGkcVGdtTZ2Z5IHyneU7zHea%2F8D93A6bcPR7e2i1XizuBW4R4oJcKPDx%2B99y5TuKD%2BU23rPlrhiW849dFG%2FstnqrVsPLJYrfYTm88mi%2BZ6uheP4DQppvcA%3D%3D&amp;_l=eNplj01rwkAYhN%2FfsleLvJv9yEbooUXoxU1BDKW5yGazNUk1SZMNtVr%2Fe9VgofQyh5mHGYYqqrhiEEYykIyDABWgAi4wlAAUKKdcRIpD3zZ%2BD5MJSBpB3dQOEI6EkdmRZGVOZgSniEySO9Kar2bwNwcvVt%2B6%2BpeJrtCQbUu7dntbmHrj1m%2FOXXOKiOfYDp3xLl%2FvTPfufLs11v1n8cKeThAni8UoCDhVAparR%2FDdsGuLxrttP7XNDlBxRjORW5pxhRYjIYULc8nDyDLm8vH3%2BRxy%2BBhM7a3p%2FdhKWSgC8XcGvuNK8%2FQp2b%2BuiiKdb1AfNNfVw6eukkBXcfk8T4L0ZVnqQ3L%2FAyjAYKo%3D&amp;_t=eNotj0tv3CAUhf1bWEcJYMeOPeomqhplMVOp8iwcVUI87tgkYCzMNPPo%2FPeAHRZwz8c5l4s0Wn5k9QFXHFeqVE9PhEAlclwRIfiBlCXGmGRy4HrMnO%2F5qGUm%2BDiCZ3LJXpEFpTnTCjVIKSFKQRUhdfGYxx1ULQtVQ81LVVNAd%2BjgvOUhel9ef0VptQUWzhNEoi3v4e9Drw%2Fx4lOrMESYYxzVALofUoo%2BLjJYw6SzkwELY%2BIRai6Ytsco0kRHyzzIwMfepGdn4F4ObJ5AouaKJu65hQB%2BTkrBPy0hVXAK6ZgnF05sxYxP0%2FeE5BYHdgoMaugdcjFLb3FlGbnHOONN1VznpmxQb5zgBm10gzcLmNwcVkBqQiPLl98mkBdVhYuIigZJHc7fJrIGPfTajYmVpI6oiiZ3HINffJTmay%2B19iop2dyy%2F93lj%2Bnandm1%2B3P3PrzvXl4%2Fu8v2tG3fzNbuL7tW6a59Hn7%2F3P%2F4Ap5ioKM%3D&amp;_b=eNozYMgoKSkottLXLy7IL6lIzkjMS0%2FVS87PZfAL9fFhsEwzsLQ0NzQ0S7GwMDRMNU8yNjA3TEpKTDM0M3MDArCqmih3V5PIrGxDX5d0g8gqr0x%2Fl0hDfxe3bN8sTwO%2FrLDMyKycDD8XEJ1uCwCmYiLM", companionAds3.CompanionClickThrough.CDATA)
						}
						if assert.NotNil(t, companionAds3.AltText) {
							assert.Equal(t, "IntegralAds_VAST_2_0_Ad_Wrapper", companionAds3.AltText)
						}

					}
				}

			}
			exts := *inline.Extensions
			if assert.Len(t, exts, 2) {
				ext1 := exts[0]
				assert.Equal(t, "LR-Pricing", ext1.Type)
				assert.Equal(t, "<Price model=\"CPM\" currency=\"USD\" source=\"spotxchange\"><![CDATA[3.06]]></Price>", strings.TrimSpace(string(ext1.Data)))
				ext2 := exts[1]
				assert.Equal(t, "SpotX-Count", ext2.Type)
				assert.Equal(t, "<total_available><![CDATA[1]]></total_available>", strings.TrimSpace(string(ext2.Data)))
			}
		}
	}
}

func TestExtraSpacesVpaid(t *testing.T) {
	v, _, _, err := loadFixture("testdata/extraspaces_vpaid.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "2.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "1130507-1818483", ad.ID)
		assert.Nil(t, ad.Wrapper)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "SpotXchange", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "IntegralAds_VAST_2_0_Ad_Wrapper", inline.AdTitle.CDATA)
			assert.Equal(t, "", inline.Description.CDATA)

			if assert.Len(t, inline.Creatives, 1) {
				crea1 := inline.Creatives[0]
				assert.Equal(t, 1, crea1.Sequence)
				if assert.NotNil(t, crea1.Linear) {
					linear := crea1.Linear

					assert.Equal(t, "        \n                  <VAST></VAST>\n                  \n                  ", linear.AdParameters.Parameters)
					if assert.Len(t, crea1.Linear.MediaFiles.MediaFile, 1) {
						media1 := crea1.Linear.MediaFiles.MediaFile
						assert.Equal(t, "progressive", media1[0].Delivery)
						assert.Equal(t, "application/javascript", media1[0].Type)
						assert.Equal(t, 300, media1[0].Width)
						assert.Equal(t, 250, media1[0].Height)
						assert.Equal(t, "\n                     https://dummy.com/dummmy.js             \n                     ", media1[0].URI)
					}
				}
			}
		}
	}
}

func TestIcons(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_adaptv_attempt_attr.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "3.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "a583680", ad.ID)
		assert.Nil(t, ad.Wrapper)
		if assert.NotNil(t, ad.InLine) {
			inline := ad.InLine
			assert.Equal(t, "Adap.tv", inline.AdSystem.Name)
			assert.Equal(t, "1.0", inline.AdSystem.Version)
			assert.Equal(t, "Adap.tv Ad Unit", inline.AdTitle.CDATA)
			//assert.Equal(t, "", inline.Description.CDATA)

			if assert.Len(t, inline.Creatives, 1) {
				crea1 := inline.Creatives[0]
				if assert.NotNil(t, crea1.Linear) {
					if assert.Len(t, *crea1.Linear.Icons.Icon, 1) {
						icon1 := *crea1.Linear.Icons.Icon
						assert.Equal(t, "DAA", icon1[0].Program)
						assert.Equal(t, 77, icon1[0].Width)
						assert.Equal(t, 15, icon1[0].Height)
						assert.Equal(t, "right", icon1[0].XPosition)
						assert.Equal(t, "top", icon1[0].YPosition)
						if assert.NotNil(t, icon1[0].StaticResource) {
							assert.Equal(t, "image/png", icon1[0].StaticResource.CreativeType)
							assert.Equal(t, "https://s.aolcdn.com/ads/adchoices.png", icon1[0].StaticResource.URI)
							assert.Equal(t, "https://adinfo.aol.com", icon1[0].IconClickThrough.CDATA)
						}
					}
				}
			}
		}
	}
}

func TestUniversalAdID(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast4_universal_ad_id.xml")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "4.0", v.Version)
	if assert.Len(t, v.Ads, 1) {
		ad := v.Ads[0]
		assert.Equal(t, "20008", ad.ID)
		if assert.NotNil(t, ad.InLine) {
			if assert.NotNil(t, ad.InLine.Extensions) {
				if assert.Len(t, ad.InLine.Creatives, 1) {
					if assert.NotNil(t, ad.InLine.Creatives[0].UniversalAdID) {
						creative := ad.InLine.Creatives[0]
						universalAdID := *creative.UniversalAdID
						assert.Equal(t, "Ad-ID", universalAdID[0].IDRegistry)
						assert.Equal(t, "8465", universalAdID[0].ID)
					}
				}
			}
		}
	}
}

func TestIABVASTSamples(t *testing.T) {
	samples := []string{
		"testdata/iab/vast_4.2_samples/Ad_Verification-test.xml",
		"testdata/iab/vast_4.2_samples/Category-test.xml",
		"testdata/iab/vast_4.2_samples/Closed_Caption_Test.xml",
		"testdata/iab/vast_4.2_samples/Event_Tracking-test.xml",
		"testdata/iab/vast_4.2_samples/IconClickFallbacks.xml",
		"testdata/iab/vast_4.2_samples/Inline_Companion_Tag-test.xml",
		"testdata/iab/vast_4.2_samples/Inline_Linear_Tag-test.xml",
		"testdata/iab/vast_4.2_samples/Inline_Non-Linear_Tag-test.xml",
		"testdata/iab/vast_4.2_samples/Inline_Simple.xml",
		"testdata/iab/vast_4.2_samples/No_Wrapper_Tag-test.xml",
		"testdata/iab/vast_4.2_samples/Ready_to_serve_Media_Files_check-test.xml",
		"testdata/iab/vast_4.2_samples/Universal_Ad_ID-multi-test.xml",
		"testdata/iab/vast_4.2_samples/Video_Clicks_and_click_tracking-Inline-test.xml",
		"testdata/iab/vast_4.2_samples/Viewable_Impression-test.xml",
		"testdata/iab/vast_4.2_samples/Wrapper_Tag-test.xml",
	}
	for _, sample := range samples {
		t.Run(sample, func(t *testing.T) {
			vast, xmlFile, _, err := loadFixture(sample)
			assert.NoError(t, err)

			expected, err := xmltree.Parse(xmlFile)
			assert.NoError(t, err)

			vastXML, err := xml.Marshal(vast)
			assert.NoError(t, err)

			actual, err := xmltree.Parse(vastXML)
			assert.NoError(t, err)

			assert.True(t, xmltree.Equal(actual, expected))
		})
	}
}
