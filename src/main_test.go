package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_messageCreate(t *testing.T) {
	type args struct {
		s *discordgo.Session
		m *discordgo.MessageCreate
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageCreate(tt.args.s, tt.args.m)
		})
	}
}
