package data

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tconnellan/macro-tracker-backend/internal/validator"
)

type Recipe struct {
	ID             int64     `json:"id"`
	Name           string    `json:"recipe_name"`
	CreatorID      int64     `json:"creator_id"`
	CreatedAt      time.Time `json:"created_at"`
	LastEditedAt   time.Time `json:"last_edited_at"`
	Notes          string    `json:"notes"`
	ParentRecipeID int64     `json:"parent_recipe_id"`
	IsLatest       bool      `json:"is_latest"`
}

func ValidateRecipe(v *validator.Validator, recipe *Recipe) {
	v.Check(recipe.Name != "", "recipe_name", "Cannot be empty")
	v.Check(len(recipe.Name) <= 50, "recipe_name", "Must be at most 50 characters")
}

// type PantryConsumable struct {
// }

type FullRecipe struct {
	Recipe           Recipe
	RecipeComponents []*RecipeComponent
	PantryItems      []*PantryItem
	Consumables      []*Consumable
}

func ValidateComponentConsumableList(v *validator.Validator, recipeID int64, recipeComponents []*RecipeComponent, pantryItems []*PantryItem, consumables []*Consumable) {
	// same length, more than zero
	v.Check(len(recipeComponents) > 0, "recipe_steps", "must have at least one step")
	v.Check(len(recipeComponents) == len(pantryItems) && len(pantryItems) == len(consumables), "recipe_steps", "must have same number of components as steps")

	if len(recipeComponents) == len(pantryItems) && len(pantryItems) == len(consumables) {
		n := int64(len(recipeComponents))
		i := int64(0)

		// sort recipeComponents based on StepNo, mirroring the alterations
		// in pantryItems and consumables. Check that step numbers are valid
		for i < n {

			for recipeComponents[i].StepNo != i+1 {

				if recipeComponents[i].StepNo < 1 || recipeComponents[i].StepNo > n {
					v.Check(false, "step_numbers", "must be in range of 1..num_steps")
					break
				}

				actualIndex := recipeComponents[i].StepNo - 1

				if recipeComponents[i].StepNo == recipeComponents[actualIndex].StepNo {
					v.Check(false, "step_numbers", "must be unique")
					break
				}

				recipeComponents[i], recipeComponents[actualIndex] = recipeComponents[actualIndex], recipeComponents[i]
				pantryItems[i], pantryItems[actualIndex] = pantryItems[actualIndex], pantryItems[i]
				consumables[i], consumables[actualIndex] = consumables[actualIndex], consumables[i]
			}

			v.Check(recipeComponents[i].PantryItemID == pantryItems[i].ID, "recipe_steps", fmt.Sprintf("pantry item ids of step %d must match", recipeComponents[i].StepNo))
			v.Check(pantryItems[i].ConsumableId == consumables[i].ID, "recipe_steps", fmt.Sprintf("consumable ids of step %d must match", recipeComponents[i].StepNo))
			v.Check(recipeComponents[i].RecipeID == recipeID, "recipe_id", "must be the same for all steps")

			i += 1
		}
	}
}

func ValidateFullRecipe(v *validator.Validator, fullRecipe *FullRecipe) {
	ValidateRecipe(v, &fullRecipe.Recipe)
	ValidateComponentConsumableList(v, fullRecipe.Recipe.ID, fullRecipe.RecipeComponents, fullRecipe.PantryItems, fullRecipe.Consumables)
	for _, recipeComponent := range fullRecipe.RecipeComponents {
		ValidateRecipeComponent(v, recipeComponent)
	}
	for _, pantryItem := range fullRecipe.PantryItems {
		ValidatePantryItem(v, pantryItem)
	}
	for _, consumable := range fullRecipe.Consumables {
		ValidateConsumable(v, consumable)
	}
}

type RecipeFilters struct {
	Metadata   MetadataFilters
	NameSearch string
}

func (r RecipeFilters) getSearchVariable() string {
	trimmedSearchVariable := strings.TrimSpace(r.NameSearch)
	if trimmedSearchVariable == "" {
		return ""
	}

	return fmt.Sprintf("%%%s%%", trimmedSearchVariable)
}

type IRecipeModel interface {
	Get(int64) (*Recipe, error)
	GetByCreatorID(int64, RecipeFilters) ([]*Recipe, Metadata, error)
	GetLatestByCreatorID(int64, RecipeFilters) ([]*Recipe, Metadata, error)
	GetFullRecipe(int64, int64) (*FullRecipe, error)
	Insert(*Recipe) error
	InsertFullRecipe(*FullRecipe) error
	Update(*Recipe) error
	UpdateFullRecipe(*FullRecipe) error
	Delete(int64) error
	GetParentRecipe(*Recipe) (*Recipe, error)
	GetAllAncestors(*Recipe, RecipeFilters) ([]*Recipe, Metadata, error)
}

type RecipeModel struct {
	DB *pgxpool.Pool
}

func (m RecipeModel) Get(ID int64) (*Recipe, error) {
	stmt := `
	SELECT id, recipe_name, creator_id, created_at, last_edited_at, notes, COALESCE(parent_recipe_id, 0), is_latest
	FROM recipes
	WHERE id = $1
	`

	var recipe Recipe

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	err := m.DB.QueryRow(ctx, stmt, ID).Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.CreatorID,
		&recipe.CreatedAt,
		&recipe.LastEditedAt,
		&recipe.Notes,
		&recipe.ParentRecipeID,
		&recipe.IsLatest,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &recipe, nil
}

func (m RecipeModel) GetByCreatorID(ID int64, filters RecipeFilters) ([]*Recipe, Metadata, error) {
	stmt := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, recipe_name, creator_id, created_at, last_edited_at, notes, COALESCE(parent_recipe_id, 0), is_latest
	FROM recipes
	WHERE creator_id = $1
	  AND ($2 = '' or recipe_name ILIKE $2)
	ORDER BY %s %s, id ASC
	LIMIT $3
	OFFSET $4
	`, filters.Metadata.sortColumn(), filters.Metadata.sortDirection())

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	rows, err := m.DB.Query(ctx, stmt, ID, filters.getSearchVariable(), filters.Metadata.pageLimit(), filters.Metadata.pageOffset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var recordCount int = 0
	var recipes []*Recipe = []*Recipe{}

	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(
			&recordCount,
			&recipe.ID,
			&recipe.Name,
			&recipe.CreatorID,
			&recipe.CreatedAt,
			&recipe.LastEditedAt,
			&recipe.Notes,
			&recipe.ParentRecipeID,
			&recipe.IsLatest,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		recipes = append(recipes, &recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return recipes, calculateMetadata(recordCount, filters.Metadata.Page, filters.Metadata.PageSize), nil
}

func (m RecipeModel) GetLatestByCreatorID(ID int64, filters RecipeFilters) ([]*Recipe, Metadata, error) {
	stmt := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, recipe_name, creator_id, created_at, last_edited_at, notes, COALESCE(parent_recipe_id, 0), is_latest
	FROM recipes
	WHERE creator_id = $1 AND is_latest = TRUE
	  AND ($2 = '' or recipe_name ILIKE $2)
	ORDER BY %s %s, id ASC
	LIMIT $3
	OFFSET $4
	`, filters.Metadata.sortColumn(), filters.Metadata.sortDirection())

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	rows, err := m.DB.Query(ctx, stmt, ID, filters.getSearchVariable(), filters.Metadata.pageLimit(), filters.Metadata.pageOffset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var recordCount int = 0
	recipes := []*Recipe{}

	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(
			&recordCount,
			&recipe.ID,
			&recipe.Name,
			&recipe.CreatorID,
			&recipe.CreatedAt,
			&recipe.LastEditedAt,
			&recipe.Notes,
			&recipe.ParentRecipeID,
			&recipe.IsLatest,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		recipes = append(recipes, &recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	return recipes, calculateMetadata(recordCount, filters.Metadata.Page, filters.Metadata.PageSize), nil
}

func (m RecipeModel) GetFullRecipe(ID int64, userID int64) (*FullRecipe, error) {
	// join recipe on componets first, then join components on consumables
	stmtRecipe := `
	SELECT id, recipe_name, creator_id, created_at, last_edited_at, notes, COALESCE(parent_recipe_id, 0), is_latest
	FROM recipes
	WHERE id = $1 AND creator_id = $2
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	txn, err := m.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.TxIsoLevel(pgx.ReadCommitted), AccessMode: pgx.ReadOnly, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		return nil, err
	}
	defer txn.Rollback(ctx)

	var recipe Recipe

	if err := txn.QueryRow(ctx, stmtRecipe, ID, userID).Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.CreatorID,
		&recipe.CreatedAt,
		&recipe.LastEditedAt,
		&recipe.Notes,
		&recipe.ParentRecipeID,
		&recipe.IsLatest,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	stmtComponents := `
	SELECT RC.id, RC.recipe_id, RC.pantry_item_id, RC.created_at, RC.quantity, RC.step_no, RC.step_description, P.id, P.user_id, P.consumable_id, P.name, P.created_at, P.last_modified, C.id, C.creator_id, C.created_at, C.name, C.brand_name, C.size, C.units, C.carbs, C.fats, C.proteins, C.alcohol
	FROM recipe_components RC 
	     INNER JOIN pantry_items P ON RC.pantry_item_id = P.id
		 INNER JOIN consumables C ON P.consumable_id = C.id
	WHERE RC.recipe_id = $1
	ORDER BY RC.step_no ASC
	`

	rows, err := txn.Query(ctx, stmtComponents, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	components := []*RecipeComponent{}
	pantryItems := []*PantryItem{}
	consumables := []*Consumable{}

	for rows.Next() {
		var component RecipeComponent
		var pantryItem PantryItem
		var consumable Consumable
		err = rows.Scan(
			&component.ID,
			&component.RecipeID,
			&component.PantryItemID,
			&component.CreatedAt,
			&component.Quantity,
			&component.StepNo,
			&component.StepDescription,
			&pantryItem.ID,
			&pantryItem.UserID,
			&pantryItem.ConsumableId,
			&pantryItem.Name,
			&pantryItem.CreatedAt,
			&pantryItem.LastEditedAt,
			&consumable.ID,
			&consumable.CreatorID,
			&consumable.CreatedAt,
			&consumable.Name,
			&consumable.BrandName,
			&consumable.Size,
			&consumable.Units,
			&consumable.Macros.Carbs,
			&consumable.Macros.Fats,
			&consumable.Macros.Proteins,
			&consumable.Macros.Alcohol,
		)
		if err != nil {
			return nil, err
		}
		components = append(components, &component)
		consumables = append(consumables, &consumable)
		pantryItems = append(pantryItems, &pantryItem)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	txn.Commit(ctx)

	return &FullRecipe{Recipe: recipe, RecipeComponents: components, PantryItems: pantryItems, Consumables: consumables}, nil

}

func (m RecipeModel) Insert(recipe *Recipe) error {
	return insertRecipe(recipe, m.DB)
}

func insertRecipe(recipe *Recipe, db psqlDB) error {
	stmt := `
	INSERT INTO recipes (recipe_name, creator_id, notes, parent_recipe_id, is_latest)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, last_edited_at
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	var actualParentID any
	if recipe.ParentRecipeID == 0 {
		actualParentID = nil
	} else {
		actualParentID = recipe.ParentRecipeID
	}

	err := db.QueryRow(ctx, stmt, recipe.Name, recipe.CreatorID, recipe.Notes, actualParentID, recipe.IsLatest).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.LastEditedAt)

	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: insert or update on table \"recipes\" violates foreign key constraint \"fk_recipe_creator\""):
			return ErrReferencedUserDoesNotExist
		case strings.HasPrefix(err.Error(), "ERROR: insert or update on table \"recipes\" violates foreign key constraint \"recipe_child_parent_id\""):
			return ErrParentRecipeDoesNotExist
		default:
			return err
		}
	}

	return nil
}

func (m RecipeModel) InsertFullRecipe(fullRecipe *FullRecipe) error {

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	// enforce ReadCommitted Isolevel and Don't defer constraint checks => guarantees consumables that
	// recipes refers to have been created
	txn, err := m.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		return err
	}
	defer txn.Rollback(ctx)

	err = insertFullRecipe(fullRecipe, txn)
	if err != nil {
		return err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return err
	}

	// making a whole other request to get the updated rows like this is not ideal
	// but doing so keeps this method in line with patterns established through the rest of our models
	// that is mutating the input model to these methods with the updated
	newFullRecipe, err := m.GetFullRecipe(fullRecipe.Recipe.ID, fullRecipe.Recipe.CreatorID)
	if err != nil {
		return err
	}

	*fullRecipe = *newFullRecipe

	return nil
}

func insertFullRecipe(fullRecipe *FullRecipe, db psqlDB) error {
	err := insertRecipe(&fullRecipe.Recipe, db)
	if err != nil {
		return err
	}

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	// insert recipe components, will fail if pantry items don't already exist
	_, err = db.CopyFrom(ctx, pgx.Identifier{"recipe_components"},
		[]string{"recipe_id", "pantry_item_id", "quantity", "step_no", "step_description"},
		pgx.CopyFromSlice(len(fullRecipe.RecipeComponents), func(i int) ([]any, error) {
			return []any{
				fullRecipe.RecipeComponents[i].RecipeID,
				fullRecipe.RecipeComponents[i].PantryItemID,
				fullRecipe.RecipeComponents[i].Quantity,
				fullRecipe.RecipeComponents[i].StepNo,
				fullRecipe.RecipeComponents[i].StepDescription,
			}, nil
		}))
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: insert or update on table \"recipe_components\" violates foreign key constraint \"fk_recipecomponent_pantry_item\""):
			return ErrPantryItemDoesNotExist
		}
		return err
	}

	return nil
}

func (m RecipeModel) Update(recipe *Recipe) error {
	return updateRecipe(recipe, m.DB)
}

func updateRecipe(recipe *Recipe, conn psqlDB) error {
	stmt := `
	UPDATE recipes
	SET recipe_name = $2, last_edited_at = current_timestamp, notes = $3, is_latest = $4
	WHERE id = $1
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	result, err := conn.Exec(ctx, stmt, recipe.ID, recipe.Name, recipe.Notes, recipe.IsLatest)
	if err != nil {
		return err
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// not good name, maybe createChildOfFullRecipe()
func (m RecipeModel) UpdateFullRecipe(fullRecipe *FullRecipe) error {

	fullRecipe.Recipe.IsLatest = false

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	txn, err := m.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		return err
	}
	defer txn.Rollback(ctx)

	err = updateRecipe(&fullRecipe.Recipe, txn)
	if err != nil {
		return err
	}

	fullRecipe.Recipe.IsLatest = true

	err = insertFullRecipe(fullRecipe, txn)
	if err != nil {
		return err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return err
	}

	// making a whole other request to get the updated rows like this is not ideal
	// but doing so keeps this method in line with patterns established through the rest of our models
	// that is mutating the input model to these methods with the updated
	newFullRecipe, err := m.GetFullRecipe(fullRecipe.Recipe.ID, fullRecipe.Recipe.CreatorID)
	if err != nil {
		return err
	}

	*fullRecipe = *newFullRecipe

	return nil
}

func (m RecipeModel) Delete(ID int64) error {

	stmtRecipeComponent := `
	DELETE FROM recipe_components 
	WHERE recipe_id = $1
	`

	stmtRecipe := `
	DELETE FROM recipes
	WHERE id = $1;
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	txn, err := m.DB.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadUncommitted, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.NotDeferrable})
	if err != nil {
		return err
	}
	defer txn.Rollback(ctx)

	_, err = txn.Exec(ctx, stmtRecipeComponent, ID)
	if err != nil {
		return err
	}

	result, err := txn.Exec(ctx, stmtRecipe, ID)
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "ERROR: update or delete on table \"recipes\" violates foreign key constraint \"recipe_child_parent_id\" on table \"recipes\""):
			return ErrChildRecipeExists
		}
		return err
	}
	rows := result.RowsAffected()

	if rows == 0 {
		return ErrRecordNotFound
	}

	err = txn.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m RecipeModel) GetParentRecipe(childRecipe *Recipe) (*Recipe, error) {
	return getParentRecipe(childRecipe, m.DB)
}

func getParentRecipe(childRecipe *Recipe, db psqlDB) (*Recipe, error) {
	if childRecipe.ParentRecipeID == 0 {
		return nil, ErrRecordNotFound
	}
	stmt := `
	SELECT id, recipe_name, creator_id, created_at, last_edited_at, notes, COALESCE(parent_recipe_id, 0), is_latest
	FROM recipes
	WHERE id = $1
	`

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	var parentRecipe Recipe

	args := []any{
		&parentRecipe.ID,
		&parentRecipe.Name,
		&parentRecipe.CreatorID,
		&parentRecipe.CreatedAt,
		&parentRecipe.LastEditedAt,
		&parentRecipe.Notes,
		&parentRecipe.ParentRecipeID,
		&parentRecipe.IsLatest,
	}

	err := db.QueryRow(ctx, stmt, childRecipe.ParentRecipeID).Scan(args...)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &parentRecipe, nil
}

func (m RecipeModel) GetAllAncestors(childRecipe *Recipe, filters RecipeFilters) ([]*Recipe, Metadata, error) {

	stmt := fmt.Sprintf(`
	WITH RECURSIVE ancestors(id, recipe_name, creator_id, created_at, last_edited_at, notes, parent_recipe_id, is_latest) AS (
		SELECT id, recipe_name, creator_id, created_at, last_edited_at, notes, COALESCE(parent_recipe_id, 0) AS parent_recipe_id, is_latest
		FROM recipes
		WHERE id = $1
		UNION
		SELECT R.id, R.recipe_name, R.creator_id, R.created_at, R.last_edited_at, R.notes, COALESCE(R.parent_recipe_id, 0) AS parent_recipe_id, R.is_latest
		FROM recipes R INNER JOIN ancestors A ON R.id = A.parent_recipe_id
	), counted_ancestors AS (
		SELECT COUNT(*) OVER() as total_count, id, recipe_name, creator_id, created_at, last_edited_at, notes, parent_recipe_id, is_latest
		FROM ancestors
	)
	SELECT total_count, id, recipe_name, creator_id, created_at, last_edited_at, notes, parent_recipe_id, is_latest
	FROM counted_ancestors
	ORDER BY %s %s, id ASC
	LIMIT $2
	OFFSET $3
	`, filters.Metadata.sortColumn(), filters.Metadata.sortDirection())

	ctx, cancel := GetDefaultTimeoutContext()
	defer cancel()

	rows, err := m.DB.Query(ctx, stmt, childRecipe.ID, filters.Metadata.pageLimit(), filters.Metadata.pageOffset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var ancestors []*Recipe
	recordCount := 0

	for rows.Next() {
		var ancestor Recipe

		err = rows.Scan(
			&recordCount,
			&ancestor.ID,
			&ancestor.Name,
			&ancestor.CreatorID,
			&ancestor.CreatedAt,
			&ancestor.LastEditedAt,
			&ancestor.Notes,
			&ancestor.ParentRecipeID,
			&ancestor.IsLatest,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		ancestors = append(ancestors, &ancestor)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	if len(ancestors) == 0 {
		return nil, Metadata{}, ErrRecordNotFound
	}

	return ancestors, calculateMetadata(recordCount, filters.Metadata.Page, filters.Metadata.PageSize), nil
}
