package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		setup    func()
		name     string
		expected *Configuration
	}{
		{
			name: "default",
			expected: &Configuration{
				KafkaBrokers: []string{"localhost:9092"},
				KafkaTopic:   "messages",
				MongoURI:     "mongodb://localhost:27017",
				MongoDB:      "message_store",
				MongoColl:    "message",
				ServerPort:   ":8080",
			},
		},
		{
			setup: func() {
				assert.NoError(t, os.Setenv("SAMPLE_KAFKA_BROKERS", "localhost:9091,localhost:9092"))
				assert.NoError(t, os.Setenv("SAMPLE_KAFKA_TOPIC", "environment"))
				assert.NoError(t, os.Setenv("SAMPLE_MONGODB_URL", "mongodb://localhost:27018"))
				assert.NoError(t, os.Setenv("SAMPLE_MONGODB_NAME", "environment"))
				assert.NoError(t, os.Setenv("SAMPLE_MONGODB_COL", "environment"))
				assert.NoError(t, os.Setenv("SAMPLE_KAFKA_TOPIC", "environment"))
				assert.NoError(t, os.Setenv("SAMPLE_SERVER_PORT", ":1234"))
			},
			name: "environment",
			expected: &Configuration{
				KafkaBrokers: []string{"localhost:9091", "localhost:9092"},
				KafkaTopic:   "environment",
				MongoURI:     "mongodb://localhost:27018",
				MongoDB:      "environment",
				MongoColl:    "environment",
				ServerPort:   ":1234",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup()
			}
			config := NewConfig()
			assert.Equal(t, test.expected, config)
		})
	}
}
