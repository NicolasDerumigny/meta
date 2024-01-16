// mautrix-meta - A Matrix-Facebook Messenger and Instagram DM puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"context"

	"go.mau.fi/util/dbutil"
	"maunium.net/go/mautrix/id"
)

const (
	getReactionByMXIDQuery = `
		SELECT message_id, thread_id, thread_receiver, reaction_sender, emoji, mxid, mx_room FROM reaction WHERE mxid=$1
	`
	getReactionByIDQuery = `
		SELECT message_id, thread_id, thread_receiver, reaction_sender, emoji, mxid, mx_room
		FROM reaction WHERE message_id=$1 AND thread_receiver=$2 AND reaction_sender=$3
	`
	insertReactionQuery = `
		INSERT INTO reaction (message_id, thread_id, thread_receiver, reaction_sender, emoji, mxid, mx_room)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	updateReactionQuery = `
		UPDATE reaction
		SET mxid=$1, emoji=$2
		WHERE message_id=$3 AND thread_receiver=$4 AND reaction_sender=$5
	`
	deleteReactionQuery = `
		DELETE FROM reaction WHERE message_id=$1 AND thread_receiver=$2 AND reaction_sender=$3
	`
)

type ReactionQuery struct {
	*dbutil.QueryHelper[*Reaction]
}

func newReaction(qh *dbutil.QueryHelper[*Reaction]) *Reaction {
	return &Reaction{qh: qh}
}

type Reaction struct {
	qh *dbutil.QueryHelper[*Reaction]

	MessageID      string
	ThreadID       int64
	ThreadReceiver int64
	Sender         int64

	Emoji string

	MXID   id.EventID
	RoomID id.RoomID
}

func (rq *ReactionQuery) GetByMXID(ctx context.Context, mxid id.EventID) (*Reaction, error) {
	return rq.QueryOne(ctx, getReactionByMXIDQuery, mxid)
}

func (rq *ReactionQuery) GetByID(ctx context.Context, msgID string, threadReceiver, reactionSender int64) (*Reaction, error) {
	return rq.QueryOne(ctx, getReactionByIDQuery, msgID, threadReceiver, reactionSender)
}

func (r *Reaction) Scan(row dbutil.Scannable) (*Reaction, error) {
	return dbutil.ValueOrErr(r, row.Scan(
		&r.MessageID, &r.ThreadID, &r.ThreadReceiver, &r.Sender, &r.Emoji, &r.MXID, &r.RoomID,
	))
}

func (r *Reaction) sqlVariables() []any {
	return []any{
		r.MessageID, r.ThreadID, r.ThreadReceiver, r.Sender, r.Emoji, r.MXID, r.RoomID,
	}
}

func (r *Reaction) Insert(ctx context.Context) error {
	return r.qh.Exec(ctx, insertReactionQuery, r.sqlVariables()...)
}

func (r *Reaction) Update(ctx context.Context) error {
	return r.qh.Exec(ctx, updateReactionQuery, r.MXID, r.Emoji, r.MessageID, r.ThreadReceiver, r.Sender)
}

func (r *Reaction) Delete(ctx context.Context) error {
	return r.qh.Exec(ctx, deleteReactionQuery, r.MessageID, r.ThreadReceiver, r.Sender)
}
