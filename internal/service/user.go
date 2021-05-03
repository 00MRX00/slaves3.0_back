package service

import (
	"fmt"

	"github.com/00mrx00/slaves3.0_back/internal/domain"
	"github.com/00mrx00/slaves3.0_back/internal/repository"
	"github.com/SevereCloud/vksdk/v2/api"
)

type AuthService struct {
	repAuth       repository.User
	repUserMaster repository.UserMaster
}

func NewAuthService(repAuth repository.User, repUserMaster repository.UserMaster) *AuthService {
	return &AuthService{
		repAuth:       repAuth,
		repUserMaster: repUserMaster,
	}
}

// func (serv *AuthService) GetUser(id int32) (domain.User, error) {
// 	user, err := serv.repAuth.GetUser(id)

// 	return user, err
// }

func (serv *AuthService) CreateUser(userId int32, userType string) (domain.UserFull, error) {
	user, err := serv.repAuth.CreateUser(userId, userType)
	if err != nil {
		return domain.UserFull{}, err
	}

	userFull, err := serv.setAddFields(user)

	return userFull, err
}

func (serv *AuthService) GetUserVkInfo(token string) (domain.UserVkInfo, error) {
	vk := api.NewVK(token)
	res, err := vk.UsersGet(api.Params{
		"fields": "screen_name, photo_100",
	})

	if err != nil {
		return domain.UserVkInfo{}, err
	}

	us := res[0]

	return domain.UserVkInfo{
		Id:        int32(us.ID),
		Firstname: us.FirstName,
		Lastname:  us.LastName,
		IsClosed:  bool(us.IsClosed),
		Username:  us.ScreenName,
		Photo:     us.Photo100,
	}, nil
}

func (serv *AuthService) GetUserFull(id int32) (domain.UserFull, error) {
	user, err := serv.repAuth.GetUser(id)
	if err != nil {
		return domain.UserFull{}, err
	}

	userFull, err := serv.setAddFields(user)

	return userFull, err
}

func (serv *AuthService) setAddFields(user domain.User) (domain.UserFull, error) {
	slaves, err := serv.repUserMaster.GetSlaves(user.Id)
	if err != nil {
		fmt.Println(err)
		return domain.UserFull{}, err
	}

	slavesCount := int32(len(slaves))

	var income int64

	for i, _ := range slaves {
		income += int64(GetSlaveProfit(slaves[i].SlaveLevel))
	}

	profit := GetSlaveProfit(user.SlaveLevel)
	damage := GetDefenderDamage(user.DefenderLevel)

	userFull := domain.UserFull{
		Id:              user.Id,
		Balance:         user.Balance,
		Gold:            user.Gold,
		SlavesCount:     slavesCount,
		Income:          income,
		Profit:          profit,
		MoneyToUpdate:   GetSlaveMoneyToUpdate(user.SlaveLevel, profit),
		Hp:              GetDefenderHp(user.DefenderLevel),
		Damage:          damage,
		DamageToUpdate:  GetDefenderDamageToUpdate(user.DefenderLevel, damage),
		LastUpdate:      user.LastUpdate,
		JobName:         user.JobName,
		UserType:        user.UserType,
		SlaveLevel:      user.SlaveLevel,
		MoneyQuantity:   user.MoneyQuantity,
		DefenderLevel:   user.DefenderLevel,
		DamageQuantity:  user.DamageQuantity,
		PurchasePriceSm: user.PurchasePriceSm,
		SalePriceSm:     GetUserSalePriceSm(user.PurchasePriceSm),
		PurchasePriceGm: user.PurchasePriceGm,
		SalePriceGm:     GetUserSalePriceGm(user.PurchasePriceGm),
		HasFetter:       GetHasFetter(user.FetterTime, user.FetterType.Duration),
		FetterTime:      user.FetterTime,
		FetterType:      user.FetterType,
		VkInfo:          user.VkInfo,
	}

	return userFull, nil
}

// func (serv *AuthService) GetFriendsList(token string, friendId int32) ([]domain.FriendInfo, error) {
// 	vk := api.NewVK(token)

// 	res, err := vk.AppsGetFriendsListExtended(api.Params{
// 		"fields": "screen_name, photo_100",
// 	})
// 	if err != nil {
// 		return []domain.FriendInfo{}, err
// 	}

// 	friends := make([]domain.FriendInfo, res.Count)

// 	for i, _ := range res.Items {
// 		frInfLoc, err := serv.GetFriendInfoLocal(friendId)
// 		if err != nil {
// 			return friends, err
// 		}

// 		if frInfLoc.MasterId != 0 {
// 			res, err := vk.UsersGet(api.Params{
// 				"fields":  "screen_name, photo_100",
// 				"user_id": frInfLoc.MasterId,
// 			})

// 			if err != nil {
// 				return friends, err
// 			}

// 			us := res[0]

// 			frInfLoc.MasterFirstname = us.FirstName
// 			frInfLoc.MasterLastname = us.LastName
// 		}

// 		friends[i] = domain.FriendInfo{
// 			Id:          int32(res.Items[i].ID),
// 			Firstname:   res.Items[i].FirstName,
// 			Lastname:    res.Items[i].LastName,
// 			Photo:       res.Items[i].Photo100,
// 			FrInfoLocal: &frInfLoc,
// 		}
// 	}

// 	return friends, nil
// }

// func (serv *AuthService) GetFriendInfoLocal(friendId int32) (domain.FriendInfoLocal, error) {
// 	frInfLoc, err := serv.repAuth.GetFriendInfoLocal(friendId)

// 	if err.Error() == "no rows in result set" {
// 		return domain.FriendInfoLocal{
// 			MasterId:        0,
// 			MasterFirstname: "",
// 			MasterLastname:  "",
// 			HasFetter:       false,
// 			FetterType:      "common",
// 			PurchasePriceSm: 20,
// 			PurchasePriceGm: 0,
// 			SlaveLevel:      0,
// 			DefenderLevel:   0,
// 		}, nil
// 	}

// 	return frInfLoc, err
// }

// func (serv *AuthService) BuySlave(userId int32, slaveId int32) error {
// 	if userId == slaveId {
// 		return errors.New("Can't buy yourself")
// 	}

// 	user, err := serv.repAuth.GetUser(userId)
// 	if err != nil {
// 		return err
// 	}

// 	slave, err := serv.repAuth.GetUser(slaveId)
// 	if err != nil {
// 		return err
// 	}

// 	timeNow, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))

// 	if slave.HasFetter {
// 		if int32(timeNow.Sub(slave.FetterTime).Minutes()) > user.FetterType.Duration {
// 			slave.HasFetter = false
// 			serv.repAuth.SetHasFetter(slave.Id, false)
// 		} else {
// 			return errors.New("UserMaster has fetter, you can't buy him")
// 		}
// 	}

// 	if user.Balance < slave.PurchasePriceSm || user.Gold < slave.PurchasePriceGm {
// 		return errors.New("Not enough money to buy a slave")
// 	}

// 	masterId, err := serv.repUserMaster.GetMaster(slaveId)
// 	if err != nil && err != pgx.ErrNoRows {
// 		return err
// 	}

// 	if masterId != 0 {
// 		if masterId == userId {
// 			return errors.New("Can't buy your slave")
// 		} else {
// 			slavesCount, balance, gold, err := serv.repAuth.GetUserBalance(masterId)
// 			if err != nil {
// 				return err
// 			}

// 			if err := serv.repAuth.SlaveCountBalanceUpdate(
// 				masterId,
// 				slavesCount-1,
// 				balance+slave.PurchasePriceSm,
// 				gold+slave.PurchasePriceGm); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	if err := serv.repAuth.SlaveCountBalanceUpdate(
// 		userId,
// 		user.SlavesCount+1,
// 		user.Balance-slave.PurchasePriceSm,
// 		user.Gold-slave.PurchasePriceGm); err != nil {
// 		return err
// 	}

// 	if err := serv.repUserMaster.CreateOrUpdateSlave(slaveId, userId); err != nil {
// 		return err
// 	}

// 	if err := serv.repAuth.SlaveBuyUpdateInfo(domain.SlaveBuyUpdateInfo{
// 		SlaveId:         slaveId,
// 		JobName:         "",
// 		UserType:        "slave",
// 		PurchasePriceSm: int64(math.Round(float64(slave.PurchasePriceSm) * 1.2)),
// 		SalePriceSm:     int64(math.Round(float64(slave.PurchasePriceSm) * 0.8)),
// 	}); err != nil {
// 		return err
// 	}

// 	return nil
// }
