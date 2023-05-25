package ark

import (
	"image"
	"reflect"
	"testing"
)

func TestScanner_ScanCharacter(t *testing.T) {
	type fields struct {
		AssetRoot string
	}
	type args struct {
		id string
	}
	f := fields{
		AssetRoot: "./test",
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CharacterArt
		wantErr bool
	}{
		{
			name:    "avg_4078_bdhkgt_1",
			fields:  f,
			args:    args{id: "avg_4078_bdhkgt_1"},
			want:    &CharacterArt{ID: "avg_4078_bdhkgt_1", Kind: "characters", BodyVariations: []CharacterArtBodyVariation{{FaceRectangle: image.Rectangle{Min: image.Point{X: 572, Y: 71}, Max: image.Point{X: 676, Y: 194}}, FaceVariations: []CharacterArtFaceVariation{{FaceSprite: "1$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "2$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "3$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "4$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "5$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "6$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "7$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "8$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "9$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "10$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "11$1.png", FaceSpriteAlpha: "1$[alpha].png", WholeBody: false}, {FaceSprite: "avg_4078_bdhkgt_1$1.png", FaceSpriteAlpha: "avg_4078_bdhkgt_1$1[alpha].png", WholeBody: false}}}, {FaceRectangle: image.Rectangle{Min: image.Point{X: 572, Y: 71}, Max: image.Point{X: 676, Y: 194}}, FaceVariations: []CharacterArtFaceVariation{{FaceSprite: "1$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "2$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "3$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "4$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "5$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "6$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "7$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "8$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "9$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "10$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "11$2.png", FaceSpriteAlpha: "2$[alpha].png", WholeBody: false}, {FaceSprite: "avg_4078_bdhkgt_1$2.png", FaceSpriteAlpha: "avg_4078_bdhkgt_1$2[alpha].png", WholeBody: false}}}, {FaceRectangle: image.Rectangle{Min: image.Point{X: 572, Y: 71}, Max: image.Point{X: 676, Y: 194}}, FaceVariations: []CharacterArtFaceVariation{{FaceSprite: "1$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "2$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "3$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "4$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "5$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "6$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "7$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "8$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "9$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "10$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "11$3.png", FaceSpriteAlpha: "3$[alpha].png", WholeBody: false}, {FaceSprite: "avg_4078_bdhkgt_1$3.png", FaceSpriteAlpha: "avg_4078_bdhkgt_1$3[alpha].png", WholeBody: false}}}, {FaceRectangle: image.Rectangle{Min: image.Point{X: 572, Y: 71}, Max: image.Point{X: 676, Y: 194}}, FaceVariations: []CharacterArtFaceVariation{{FaceSprite: "1$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "2$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "3$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "4$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "5$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "6$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "7$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "8$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "9$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "10$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "11$4.png", FaceSpriteAlpha: "4$[alpha].png", WholeBody: false}, {FaceSprite: "avg_4078_bdhkgt_1$4.png", FaceSpriteAlpha: "avg_4078_bdhkgt_1$4[alpha].png", WholeBody: false}}}}},
			wantErr: false,
		},
		{
			name:    "avg_npc_034",
			fields:  f,
			args:    args{id: "avg_npc_034"},
			want:    &CharacterArt{ID: "avg_npc_034", Kind: "characters", BodyVariations: []CharacterArtBodyVariation{{FaceRectangle: image.Rectangle{Min: image.Point{X: -1, Y: -1}, Max: image.Point{X: -1, Y: -1}}, FaceVariations: []CharacterArtFaceVariation{{FaceSprite: "avg_npc_034.png", FaceSpriteAlpha: "avg_npc_034[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_2.png", FaceSpriteAlpha: "avg_npc_034_2[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_3.png", FaceSpriteAlpha: "avg_npc_034_3[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_4.png", FaceSpriteAlpha: "avg_npc_034_4[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_5.png", FaceSpriteAlpha: "avg_npc_034_5[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_6.png", FaceSpriteAlpha: "avg_npc_034_6[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_7.png", FaceSpriteAlpha: "avg_npc_034_7[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_8.png", FaceSpriteAlpha: "avg_npc_034_8[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_9.png", FaceSpriteAlpha: "avg_npc_034_9[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_10.png", FaceSpriteAlpha: "avg_npc_034_10[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_11.png", FaceSpriteAlpha: "avg_npc_034_11[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_12.png", FaceSpriteAlpha: "avg_npc_034_12[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_13.png", FaceSpriteAlpha: "avg_npc_034_13[alpha].png", WholeBody: false}, {FaceSprite: "avg_npc_034_14.png", FaceSpriteAlpha: "avg_npc_034_14[alpha].png", WholeBody: false}}}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := &Scanner{
				Root: tt.fields.AssetRoot,
			}
			got, err := scanner.ScanCharacter(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanCharacter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanCharacter() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
