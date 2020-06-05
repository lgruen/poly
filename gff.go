package main

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
)

// ParseGff Takes in a string representing a gffv3 file and parses it into an AnnotatedSequence object.
func ParseGff(gff string) AnnotatedSequence {
	lines := strings.Split(gff, "\n")
	metaString := lines[0:2]
	versionString := metaString[0]
	regionStringArray := strings.Split(metaString[1], " ")

	meta := Meta{}
	meta.GffVersion = strings.Split(versionString, " ")[1]
	meta.Name = regionStringArray[1] // Formally region name, but changed to name here for generality/interoperability.
	meta.RegionStart, _ = strconv.Atoi(regionStringArray[2])
	meta.RegionEnd, _ = strconv.Atoi(regionStringArray[3])
	meta.Size = meta.RegionEnd - meta.RegionStart

	records := []Feature{}
	sequence := Sequence{}
	var sequenceBuffer bytes.Buffer
	fastaFlag := false
	for _, line := range lines {
		if line == "##FASTA" {
			fastaFlag = true
		} else if len(line) == 0 {
			continue
		} else if line[0:2] == "##" {
			continue
		} else if fastaFlag == true && line[0:1] != ">" {
			// sequence.Sequence = sequence.Sequence + line
			sequenceBuffer.WriteString(line)
		} else if fastaFlag == true && line[0:1] == ">" {
			sequence.Description = line
		} else {
			record := Feature{}
			fields := strings.Split(line, "\t")
			record.Name = fields[0]
			record.Source = fields[1]
			record.Type = fields[2]
			record.Start, _ = strconv.Atoi(fields[3])
			record.End, _ = strconv.Atoi(fields[4])
			record.Score, _ = strconv.ParseFloat(fields[5], 64)
			record.Strand = fields[6][0]
			record.Phase, _ = strconv.Atoi(fields[7])
			record.Attributes = make(map[string]string)
			attributes := fields[8]
			var eqIndex int
			for i := strings.Index(attributes, ";"); i > 0; i = strings.Index(attributes, ";") {
				eqIndex = strings.Index(attributes[:i], "=")
				record.Attributes[attributes[:i][:eqIndex]] = attributes[:i][eqIndex+1:]
				attributes = attributes[i+1:]
			}
			records = append(records, record)
		}
	}
	sequence.Sequence = sequenceBuffer.String()
	annotatedSequence := AnnotatedSequence{}
	annotatedSequence.Meta = meta
	annotatedSequence.Features = records
	annotatedSequence.Sequence = sequence

	return annotatedSequence
}

// BuildGff takes an Annotated sequence and returns a string representing a gff to be written out.
// func BuildGff(AnnotatedSequence) string {
// 	var Gff string

// }

// ReadGff takes in a filepath for a .gffv3 file and parses it into an Annotated Sequence struct.
func ReadGff(path string) AnnotatedSequence {
	file, err := ioutil.ReadFile(path)
	var annotatedSequence AnnotatedSequence
	if err != nil {
		// return 0, fmt.Errorf("Failed to open file %s for unpack: %s", gzFilePath, err)
	} else {
		annotatedSequence = ParseGff(string(file))
	}
	return annotatedSequence
}

// WriteGff takes an AnnotatedSequence struct and a path string and writes out a gff to that path.
// func WriteGff(annotatedSequence AnnotatedSequence, path string) AnnotatedSequence {
// 	file
// 	_ = ioutil.WriteFile(path, file, 0644)

// }
