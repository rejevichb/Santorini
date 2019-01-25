package board

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

// This is very bad, we should be using [][]Cell instead of this monstrosity
func (b board) MarshalJSON() ([]byte, error) {
	width, height := b.Dimensions()
	buffer := bytes.NewBufferString("[")
	for x := 0; x < width; x++ {
		buffer.WriteString("[")

		for y := 0; y < height; y++ {
			pos := Pos{y, x}
			ITile, err := b.TileAt(pos)
			if err != nil {
				return nil, err
			}
			buffer.WriteString("\"" + strconv.Itoa(ITile.FloorCount()))

			worker := b.WorkerAt(pos)
			if worker != nil {
				buffer.WriteString(worker.Name())
			}

			buffer.WriteString("\"")

			if y < height-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("]")

		if x < width-1 {
			buffer.WriteString(",\n")
		}
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

//Unmarshals the JSON Array into a board
func (b *board) UnmarshalJSON(buf []byte) error {
	var cells [][]Cell
	if err := json.Unmarshal(buf, &cells); err != nil {
		return err
	}

	workers := make([]IWorker, 0)

	for y, row := range cells {
		for x, cell := range row {
			p := Pos{X: x, Y: y}
			b.tiles[x][y] = CustomTile(p, cell.Height)

			if cell.Worker != "" {
				// find the ID of the worker being touched and add it to an outer
				// worker array, to add at the end
				// NOTE save the name, ID, and Pos of the worker
				player, workerID, err := ParseWorkerName(cell.Worker)
				if err != nil {
					return err
				}
				worker := NewWorker(p, player, workerID)
				workers = append(workers, worker)
			}
		}
	}

	copy(b.workers[:len(b.workers)], workers[:len(b.workers)])

	return nil
}

type Cell struct {
	Height int
	Worker string
}

//Unmarshals the JSON interface into a struct
func (c *Cell) UnmarshalJSON(buf []byte) error {
	var tmp string
	if err := json.Unmarshal(buf, &tmp); err == nil {
		c.Height, err = strconv.Atoi(tmp[:1])
		if err != nil {
			return errors.New("Unable to acquire height of BuildingWorker")
		}
		c.Worker = tmp[1:]
	} else {
		if err := json.Unmarshal(buf, &c.Height); err != nil {
			return err
		}
	}
	return nil
}
