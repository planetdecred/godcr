package wallet

import (
	"github.com/planetdecred/dcrlibwallet"
)

type AgendaStatus int

const (
	SyncedAgenda AgendaStatus = iota
	AgendaInProgress
	NewAgendaFound
	AgendaFinished
)

type Agenda struct {
	Agenda       *dcrlibwallet.Agenda
	AgendaStatus AgendaStatus
}

func (l *listener) OnNewAgenda(agenda *dcrlibwallet.Agenda) {
	l.Send <- SyncStatusUpdate{
		Stage: AgendaAdded,
		Agenda: Agenda{
			Agenda:       agenda,
			AgendaStatus: NewAgendaFound,
		},
	}
}

func (l *listener) OnAgendaVoteStarted(agenda *dcrlibwallet.Agenda) {
	l.Send <- SyncStatusUpdate{
		Stage: AgendaVoteStarted,
		Agenda: Agenda{
			Agenda:       agenda,
			AgendaStatus: AgendaInProgress,
		},
	}
}

func (l *listener) OnAgendaVoteFinished(agenda *dcrlibwallet.Agenda) {
	l.Send <- SyncStatusUpdate{
		Stage: AgendaVoteFinished,
		Agenda: Agenda{
			Agenda:       agenda,
			AgendaStatus: AgendaFinished,
		},
	}
}

func (l *listener) OnAgendasSynced() {
	l.Send <- SyncStatusUpdate{
		Stage: AgendaSynced,
		Agenda: Agenda{
			AgendaStatus: SyncedAgenda,
		},
	}
}
