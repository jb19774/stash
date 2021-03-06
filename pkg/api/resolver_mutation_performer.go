package api

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) PerformerCreate(ctx context.Context, input models.PerformerCreateInput) (*models.Performer, error) {
	// generate checksum from performer name rather than image
	checksum := utils.MD5FromString(input.Name)

	var imageData []byte
	var err error

	if input.Image != nil {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)
	}

	if err != nil {
		return nil, err
	}

	// Populate a new performer from the input
	currentTime := time.Now()
	newPerformer := models.Performer{
		Checksum:  checksum,
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}
	newPerformer.Name = sql.NullString{String: input.Name, Valid: true}
	if input.URL != nil {
		newPerformer.URL = sql.NullString{String: *input.URL, Valid: true}
	}
	if input.Gender != nil {
		newPerformer.Gender = sql.NullString{String: input.Gender.String(), Valid: true}
	}
	if input.Birthdate != nil {
		newPerformer.Birthdate = models.SQLiteDate{String: *input.Birthdate, Valid: true}
	}
	if input.Ethnicity != nil {
		newPerformer.Ethnicity = sql.NullString{String: *input.Ethnicity, Valid: true}
	}
	if input.Country != nil {
		newPerformer.Country = sql.NullString{String: *input.Country, Valid: true}
	}
	if input.EyeColor != nil {
		newPerformer.EyeColor = sql.NullString{String: *input.EyeColor, Valid: true}
	}
	if input.Height != nil {
		newPerformer.Height = sql.NullString{String: *input.Height, Valid: true}
	}
	if input.Measurements != nil {
		newPerformer.Measurements = sql.NullString{String: *input.Measurements, Valid: true}
	}
	if input.FakeTits != nil {
		newPerformer.FakeTits = sql.NullString{String: *input.FakeTits, Valid: true}
	}
	if input.CareerLength != nil {
		newPerformer.CareerLength = sql.NullString{String: *input.CareerLength, Valid: true}
	}
	if input.Tattoos != nil {
		newPerformer.Tattoos = sql.NullString{String: *input.Tattoos, Valid: true}
	}
	if input.Piercings != nil {
		newPerformer.Piercings = sql.NullString{String: *input.Piercings, Valid: true}
	}
	if input.Aliases != nil {
		newPerformer.Aliases = sql.NullString{String: *input.Aliases, Valid: true}
	}
	if input.Twitter != nil {
		newPerformer.Twitter = sql.NullString{String: *input.Twitter, Valid: true}
	}
	if input.Instagram != nil {
		newPerformer.Instagram = sql.NullString{String: *input.Instagram, Valid: true}
	}
	if input.Favorite != nil {
		newPerformer.Favorite = sql.NullBool{Bool: *input.Favorite, Valid: true}
	} else {
		newPerformer.Favorite = sql.NullBool{Bool: false, Valid: true}
	}

	// Start the transaction and save the performer
	var performer *models.Performer
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()

		performer, err = qb.Create(newPerformer)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(performer.ID, imageData); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if input.StashIds != nil {
			stashIDJoins := models.StashIDsFromInput(input.StashIds)
			if err := qb.UpdateStashIDs(performer.ID, stashIDJoins); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return performer, nil
}

func (r *mutationResolver) PerformerUpdate(ctx context.Context, input models.PerformerUpdateInput) (*models.Performer, error) {
	// Populate performer from the input
	performerID, _ := strconv.Atoi(input.ID)
	updatedPerformer := models.PerformerPartial{
		ID:        performerID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: time.Now()},
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	var imageData []byte
	var err error
	imageIncluded := translator.hasField("image")
	if input.Image != nil {
		_, imageData, err = utils.ProcessBase64Image(*input.Image)
		if err != nil {
			return nil, err
		}
	}

	if input.Name != nil {
		// generate checksum from performer name rather than image
		checksum := utils.MD5FromString(*input.Name)

		updatedPerformer.Name = &sql.NullString{String: *input.Name, Valid: true}
		updatedPerformer.Checksum = &checksum
	}

	updatedPerformer.URL = translator.nullString(input.URL, "url")

	if translator.hasField("gender") {
		if input.Gender != nil {
			updatedPerformer.Gender = &sql.NullString{String: input.Gender.String(), Valid: true}
		} else {
			updatedPerformer.Gender = &sql.NullString{String: "", Valid: false}
		}
	}

	updatedPerformer.Birthdate = translator.sqliteDate(input.Birthdate, "birthdate")
	updatedPerformer.Country = translator.nullString(input.Country, "country")
	updatedPerformer.EyeColor = translator.nullString(input.EyeColor, "eye_color")
	updatedPerformer.Measurements = translator.nullString(input.Measurements, "measurements")
	updatedPerformer.Height = translator.nullString(input.Height, "height")
	updatedPerformer.Ethnicity = translator.nullString(input.Ethnicity, "ethnicity")
	updatedPerformer.FakeTits = translator.nullString(input.FakeTits, "fake_tits")
	updatedPerformer.CareerLength = translator.nullString(input.CareerLength, "career_length")
	updatedPerformer.Tattoos = translator.nullString(input.Tattoos, "tattoos")
	updatedPerformer.Piercings = translator.nullString(input.Piercings, "piercings")
	updatedPerformer.Aliases = translator.nullString(input.Aliases, "aliases")
	updatedPerformer.Twitter = translator.nullString(input.Twitter, "twitter")
	updatedPerformer.Instagram = translator.nullString(input.Instagram, "instagram")
	updatedPerformer.Favorite = translator.nullBool(input.Favorite, "favorite")

	// Start the transaction and save the performer
	var performer *models.Performer
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()

		var err error
		performer, err = qb.Update(updatedPerformer)
		if err != nil {
			return err
		}

		// update image table
		if len(imageData) > 0 {
			if err := qb.UpdateImage(performer.ID, imageData); err != nil {
				return err
			}
		} else if imageIncluded {
			// must be unsetting
			if err := qb.DestroyImage(performer.ID); err != nil {
				return err
			}
		}

		// Save the stash_ids
		if translator.hasField("stash_ids") {
			stashIDJoins := models.StashIDsFromInput(input.StashIds)
			if err := qb.UpdateStashIDs(performerID, stashIDJoins); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return performer, nil
}

func (r *mutationResolver) PerformerDestroy(ctx context.Context, input models.PerformerDestroyInput) (bool, error) {
	id, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		return repo.Performer().Destroy(id)
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) PerformersDestroy(ctx context.Context, performerIDs []string) (bool, error) {
	ids, err := utils.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Performer()
		for _, id := range ids {
			if err := qb.Destroy(id); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}
