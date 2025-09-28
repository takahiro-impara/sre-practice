package postgres

import (
	db "github.com/lot-koichi/sre-skill-up-project/services/user/db/sqlc/generated"
	"github.com/lot-koichi/sre-skill-up-project/services/user/internal/domain"
)

// toDomainUser converts SQLC generated User to domain User
func toDomainUser(sqlcUser db.User) *domain.User {
	return &domain.User{
		ID:        sqlcUser.ID,
		Email:     domain.Email(sqlcUser.Email),
		Password:  domain.Password(sqlcUser.Password),
		Name:      domain.Name(sqlcUser.Name),
		CreatedAt: sqlcUser.CreatedAt.Time,
		UpdatedAt: sqlcUser.UpdatedAt.Time,
	}
}

// toCreateUserParams converts domain User to SQLC CreateUserParams
func toCreateUserParams(user *domain.User) db.CreateUserParams {
	return db.CreateUserParams{
		Email:    string(user.Email),
		Password: string(user.Password),
		Name:     string(user.Name),
	}
}

// toUpdateUserParams converts domain User to SQLC UpdateUserParams
func toUpdateUserParams(user *domain.User) db.UpdateUserParams {
	return db.UpdateUserParams{
		Email: string(user.Email),
		Name:  string(user.Name),
		ID:    user.ID,
	}
}

// toListUsersParams creates SQLC ListUsersParams
func toListUsersParams(limit, offset int32) db.ListUsersParams {
	return db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}
}

// toDomainUsers converts multiple SQLC Users to domain Users
func toDomainUsers(sqlcUsers []db.User) []*domain.User {
	domainUsers := make([]*domain.User, 0, len(sqlcUsers))
	for _, sqlcUser := range sqlcUsers {
		domainUsers = append(domainUsers, toDomainUser(sqlcUser))
	}
	return domainUsers
}

// updateTimestamps updates the domain user with timestamps from created/updated user
func updateTimestamps(domainUser *domain.User, sqlcUser db.User) {
	domainUser.ID = sqlcUser.ID
	domainUser.CreatedAt = sqlcUser.CreatedAt.Time
	domainUser.UpdatedAt = sqlcUser.UpdatedAt.Time
}

// updateUpdatedAt updates only the UpdatedAt timestamp
func updateUpdatedAt(domainUser *domain.User, sqlcUser db.User) {
	domainUser.UpdatedAt = sqlcUser.UpdatedAt.Time
}

