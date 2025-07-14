package srlogrotate

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLogger_fileName(t *testing.T) {
	t.Run("日付の変わりめ1", func(t *testing.T) {
		loc, err := time.LoadLocation("Asia/Tokyo")
		assert.NoError(t, err)

		l := logger{
			fileBaseName: "log_file",
			timeFormat:   "20060102",
			nowFunc: func() time.Time {
				return time.Date(2024, 2, 3, 0, 0, 0, 0, loc)
			},
		}

		actual := l.fileName()
		expected := "log_file.20240203"
		assert.Equal(t, expected, actual)
	})

	t.Run("日付の変わりめ1", func(t *testing.T) {
		loc, err := time.LoadLocation("Asia/Tokyo")
		assert.NoError(t, err)

		l := logger{
			fileBaseName: "log_file",
			timeFormat:   "20060102",
			nowFunc: func() time.Time {
				return time.Date(2024, 2, 3, 0, 0, 1, 0, loc)
			},
		}

		actual := l.fileName()
		expected := "log_file.20240203"
		assert.Equal(t, expected, actual)
	})
}

func TestLogger_Write(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	assert.NoError(t, err)

	t.Run("test", func(t *testing.T) {
		l := &logger{
			fileBaseName: "log_file",
			timeFormat:   "20060102",
			nowFunc: func() time.Time {
				return time.Date(2024, 2, 2, 23, 59, 55, 0, loc)
			},
		}

		var err error
		_, err = l.Write([]byte("aaaa\n"))
		assert.NoError(t, err)

		l.nowFunc = func() time.Time {
			return time.Date(2024, 2, 2, 23, 59, 59, 0, loc)
		}
		_, err = l.Write([]byte("bbbb\n"))
		assert.NoError(t, err)

		l.nowFunc = func() time.Time {
			return time.Date(2024, 2, 3, 0, 0, 0, 0, loc)
		}
		_, err = l.Write([]byte("cccc\n"))
		assert.NoError(t, err)
	})
}
