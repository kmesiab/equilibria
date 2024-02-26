package models

type MessageStatus struct {
	ID   int64  `json:"status_id" gorm:"primaryKey;autoIncrement"`
	Name string ` json:"status_name" gorm:"size:100;unique;not null"`
}

func NewMessageStatusPending() MessageStatus {
	return MessageStatus{
		ID:   1,
		Name: "Pending Activation",
	}
}

func NewMessageStatusSent() MessageStatus {
	return MessageStatus{
		ID:   2,
		Name: "Sent",
	}
}

func NewMessageStatusReceived() MessageStatus {
	return MessageStatus{
		ID:   3,
		Name: "Received",
	}
}

func NewMessageStatusDelivered() MessageStatus {
	return MessageStatus{
		ID:   4,
		Name: "Delivered",
	}
}

func NewMessageStatusCanceled() MessageStatus {
	return MessageStatus{
		ID:   5,
		Name: "Canceled",
	}
}

func NewMessageStatusFailed() MessageStatus {
	return MessageStatus{
		ID:   6,
		Name: "Failed",
	}
}

func NewMessageStatusAccepted() MessageStatus {
	return MessageStatus{
		ID:   7,
		Name: "Accepted",
	}
}

func NewMessageStatusQueued() MessageStatus {
	return MessageStatus{
		ID:   8,
		Name: "Queued",
	}
}

func NewMessageStatusReceiving() MessageStatus {
	return MessageStatus{
		ID:   9,
		Name: "Receiving",
	}
}

func NewMessageStatusRead() MessageStatus {
	return MessageStatus{
		ID:   10,
		Name: "Read",
	}
}

func NewMessageStatusSending() MessageStatus {
	return MessageStatus{
		ID:   11,
		Name: "Sending",
	}
}

func NewMessageStatusUnknown() MessageStatus {
	return MessageStatus{
		ID:   12,
		Name: "Unknown",
	}
}
