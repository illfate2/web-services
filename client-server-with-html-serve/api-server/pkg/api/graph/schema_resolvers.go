package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/api/generated"
	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/api/generated/model"
	"github.com/illfate2/web-services/client-server-with-html-serve/api-server/pkg/entities"
)

func (r *mutationResolver) CreateMuseumItem(ctx context.Context, input model.MuseumItemInput) (*model.MuseumItem, error) {
	museumItem, err := r.service.CreateMuseumItem(entities.MuseumItemWithDetails{
		MuseumItem: entities.MuseumItem{
			InventoryNumber: input.InventoryNumber,
			Name:            input.Name,
			CreationDate: entities.Date{
				Time: input.CreationDate,
			},
			MuseumSetID:  input.SetID,
			MuseumFundID: input.FundID,
			Annotation:   getStrFromPtr(input.Annotation),
		},
		Keeper: entities.Person{
			FirstName:  input.PersonInput.FirstName,
			LastName:   input.PersonInput.LastName,
			MiddleName: input.PersonInput.MiddleName,
		},
	})
	if err != nil {
		return nil, err
	}
	return convertEntityMuseumItem(museumItem.MuseumItem), nil
}

func convertEntityMuseumItem(item entities.MuseumItem) *model.MuseumItem {
	return &model.MuseumItem{
		ID:              item.ID,
		InventoryNumber: item.InventoryNumber,
		Name:            item.Name,
		CreationDate:    item.CreationDate.Time,
		Annotation:      &item.Annotation,
		Set: &model.MuseumSet{
			ID: item.MuseumSetID,
		},
		Fund: &model.MuseumFund{
			ID: item.MuseumFundID,
		},
		Person: &model.Person{
			ID: &item.KeeperID,
		},
	}
}

func (r *mutationResolver) UpdateMuseumItem(ctx context.Context, id int, input model.UpdateMuseumItemInput) (*model.MuseumItem, error) {
	err := r.service.UpdateMuseumItem(entities.MuseumItem{
		ID:              id,
		InventoryNumber: *input.InventoryNumber,
		Name:            *input.Name,
		CreationDate:    entities.NewDate(*input.CreationDate),
		Annotation:      getStrFromPtr(input.Annotation),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update museum item")
	}
	return r.Resolver.Query().MuseumItem(ctx, id)
}

func (r *mutationResolver) DeleteMuseumItem(ctx context.Context, id int) (int, error) {
	err := r.service.DeleteMuseumItem(id)
	return id, err
}

func (r *mutationResolver) CreateMuseumSet(ctx context.Context, input model.MuseumSetInput) (*model.MuseumSet, error) {
	set, err := r.service.CreateMuseumSet(entities.MuseumSet{
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}
	return convertEntityMuseumSet(set), nil
}

func convertEntityMuseumSet(set entities.MuseumSet) *model.MuseumSet {
	return &model.MuseumSet{
		ID:   set.ID,
		Name: set.Name,
	}
}

func (r *mutationResolver) UpdateMuseumSet(ctx context.Context, id int, input model.MuseumSetInput) (*model.MuseumSet, error) {
	err := r.service.UpdateMuseumSet(entities.MuseumSet{
		ID:   id,
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}
	return r.Resolver.Query().MuseumSet(ctx, id)
}

func (r *mutationResolver) DeleteMuseumSet(ctx context.Context, id int) (int, error) {
	err := r.service.DeleteMuseumSet(id)
	return id, err
}

func (r *mutationResolver) CreateMuseumItemMovement(ctx context.Context, input model.MuseumMovementInput) (*model.MuseumItemMovement, error) {
	movement, err := r.service.CreateMuseumItemMovement(entities.MuseumItemMovement{
		MuseumItemID:        input.ItemID,
		AcceptDate:          input.AcceptDate,
		ExhibitTransferDate: input.ExhibitTransferDate,
		ExhibitReturnDate:   input.ExhibitReturnDate,
		ResponsiblePerson: entities.Person{
			FirstName:  input.PersonInput.FirstName,
			LastName:   input.PersonInput.LastName,
			MiddleName: input.PersonInput.MiddleName,
		},
	})
	if err != nil {
		return nil, err
	}
	return convertEntityMuseumMovement(movement), nil
}

func convertEntityMuseumMovement(m entities.MuseumItemMovement) *model.MuseumItemMovement {
	return &model.MuseumItemMovement{
		ID:                  m.ID,
		AcceptDate:          m.AcceptDate,
		ExhibitTransferDate: m.ExhibitTransferDate,
		ExhibitReturnDate:   m.ExhibitReturnDate,
		Person: &model.Person{
			ID: &m.ResponsiblePersonID,
		},
		Item: &model.MuseumItem{
			ID: m.MuseumItemID,
		},
	}
}

func (r *mutationResolver) UpdateMuseumItemMovement(ctx context.Context, id int, input model.UpdateMuseumMovementInput) (*model.MuseumItemMovement, error) {
	err := r.service.UpdateMuseumItemMovement(entities.MuseumItemMovement{
		ID:                  id,
		AcceptDate:          input.AcceptDate,
		ExhibitTransferDate: input.ExhibitTransferDate,
		ExhibitReturnDate:   input.ExhibitReturnDate,
	})
	if err != nil {
		return nil, err
	}
	return r.Resolver.Query().MuseumMovement(ctx, id)
}

func (r *mutationResolver) DeleteMuseumItemMovement(ctx context.Context, id int) (int, error) {
	err := r.service.DeleteMuseumItemMovement(id)
	return id, err
}

func (r *mutationResolver) CreateMuseumFund(ctx context.Context, input model.MuseumFundInput) (*model.MuseumFund, error) {
	fund, err := r.service.CreateMuseumFund(entities.MuseumFund{
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}
	return convertEntityFund(fund), nil
}

func convertEntityFund(fund entities.MuseumFund) *model.MuseumFund {
	return &model.MuseumFund{
		ID:   fund.ID,
		Name: fund.Name,
	}
}

func (r *mutationResolver) UpdateMuseumFund(ctx context.Context, id int, input model.UpdateMuseumFundInput) (*model.MuseumFund, error) {
	err := r.service.UpdateMuseumFund(entities.MuseumFund{
		ID:   id,
		Name: *input.Name,
	})
	if err != nil {
		return nil, err
	}
	return r.Resolver.Query().MuseumFund(ctx, id)
}

func (r *mutationResolver) DeleteMuseumFund(ctx context.Context, id int) (int, error) {
	err := r.service.DeleteMuseumFund(id)
	return id, err
}

func (r *queryResolver) MuseumItem(ctx context.Context, id int) (*model.MuseumItem, error) {
	item, err := r.service.GetMuseumItem(id)
	if err != nil {
		return nil, err
	}
	return convertEntityMuseumItem(item), nil
}

func (r *queryResolver) MuseumItems(ctx context.Context) ([]*model.MuseumItem, error) {
	items, err := r.service.SearchMuseumItems(entities.SearchMuseumItemsArgs{})
	if err != nil {
		return nil, err
	}
	res := make([]*model.MuseumItem, 0, len(items))
	for _, i := range items {
		res = append(res, convertEntityMuseumItem(i))
	}
	return res, nil
}

func (r *queryResolver) MuseumSet(ctx context.Context, id int) (*model.MuseumSet, error) {
	set, err := r.service.GetMuseumSet(id)
	if err != nil {
		log.Print("faild to get set", err)
		return nil, err
	}
	return convertEntityMuseumSet(set.MuseumSet), nil
}

func (r *queryResolver) MuseumSets(ctx context.Context) ([]*model.MuseumSet, error) {
	sets, err := r.service.GetMuseumSets()
	if err != nil {
		return nil, err
	}
	modelsSets := make([]*model.MuseumSet, 0, len(sets))
	for _, s := range sets {
		modelsSets = append(modelsSets, convertEntityMuseumSet(s))
	}
	return modelsSets, nil
}

func (r *queryResolver) MuseumFund(ctx context.Context, id int) (*model.MuseumFund, error) {
	fund, err := r.service.GetMuseumFundByID(id)
	if err != nil {
		log.Print("failed to get fund", err)
		return nil, err
	}
	return convertEntityFund(fund), nil
}

func (r *queryResolver) MuseumFunds(ctx context.Context) ([]*model.MuseumFund, error) {
	funds, err := r.service.GetMuseumFunds()
	if err != nil {
		return nil, err
	}
	res := make([]*model.MuseumFund, 0, len(funds))
	for _, f := range funds {
		res = append(res, convertEntityFund(f))
	}
	return res, nil
}

func (r *queryResolver) MuseumMovement(ctx context.Context, id int) (*model.MuseumItemMovement, error) {
	movement, err := r.service.GetMuseumItemMovement(id)
	if err != nil {
		return nil, err
	}
	return convertEntityMuseumMovement(movement), nil
}

func (r *queryResolver) MuseumMovements(ctx context.Context) ([]*model.MuseumItemMovement, error) {
	movements, err := r.service.GetMuseumItemMovements()
	if err != nil {
		return nil, err
	}
	res := make([]*model.MuseumItemMovement, 0, len(movements))
	for _, m := range movements {
		res = append(res, convertEntityMuseumMovement(m))
	}
	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *Resolver) MuseumItem() generated.MuseumItemResolver {
	return &museumItemResolver{r}
}

func (r *Resolver) MuseumItemMovement() generated.MuseumItemMovementResolver {
	return &museumItemMovementResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type museumItemResolver struct {
	*Resolver
}

func (r *museumItemResolver) Person(ctx context.Context, obj *model.MuseumItem) (*model.Person, error) {
	p, err := r.service.FindPerson(*obj.Person.ID)
	if err != nil {
		return nil, err
	}
	return &model.Person{
		ID:         &p.ID,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		MiddleName: p.MiddleName,
	}, nil
}

func (r *museumItemResolver) Fund(ctx context.Context, item *model.MuseumItem) (*model.MuseumFund, error) {
	return r.Query().MuseumFund(ctx, item.Fund.ID)
}

func (r *museumItemResolver) Set(ctx context.Context, obj *model.MuseumItem) (*model.MuseumSet, error) {
	return r.Query().MuseumSet(ctx, obj.Set.ID)
}

func getStrFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type museumItemMovementResolver struct {
	*Resolver
}

func (r *museumItemMovementResolver) Item(ctx context.Context, m *model.MuseumItemMovement) (*model.MuseumItem, error) {
	return r.Query().MuseumItem(ctx, m.Item.ID)
}

func (r *museumItemMovementResolver) Person(ctx context.Context, m *model.MuseumItemMovement) (*model.Person, error) {
	p, err := r.service.FindPerson(*m.Person.ID)
	if err != nil {
		return nil, err
	}
	return &model.Person{
		ID:         &p.ID,
		FirstName:  p.FirstName,
		LastName:   p.LastName,
		MiddleName: p.MiddleName,
	}, nil
}
