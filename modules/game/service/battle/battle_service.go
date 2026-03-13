package battle

import (
	"context"
	"gamebook-backend/database/entities"
	battleCommon "gamebook-backend/modules/game/battle"
	"gamebook-backend/modules/game/dto"
	gameDTO "gamebook-backend/modules/game/dto"
	"gamebook-backend/modules/game/log"
	"gamebook-backend/modules/game/repository"

	"gorm.io/gorm"
)

type Service interface {
	Action(ctx context.Context, player entities.Player, req gameDTO.ActionRequest) (gameDTO.ActionResponse, error)
	Battle(ctx context.Context, player entities.Player, req dto.ActionRequest) (dto.BattleDto, error)
}

type service struct {
	sectionRepository            repository.SectionRepository
	dicesRepository              repository.DiceRepository
	battleRepository             repository.BattleRepository
	playerRepository             repository.PlayerRepository
	playerSectionEnemyRepository repository.PlayerSectionEnemyRepository
	db                           *gorm.DB
	logService                   log.PlayerLogService
}

func NewService(
	sectionRepo repository.SectionRepository,
	dicesRepo repository.DiceRepository,
	battleRepo repository.BattleRepository,
	playerRepo repository.PlayerRepository,
	db *gorm.DB,
) Service {
	return &service{
		sectionRepository: sectionRepo,
		dicesRepository:   dicesRepo,
		battleRepository:  battleRepo,
		playerRepository:  playerRepo,
		db:                db,
	}
}

func NewServiceWithLogging(
	sectionRepo repository.SectionRepository,
	dicesRepo repository.DiceRepository,
	battleRepo repository.BattleRepository,
	playerRepo repository.PlayerRepository,
	playerSectionEnemyRepo repository.PlayerSectionEnemyRepository,
	db *gorm.DB,
	logService log.PlayerLogService,
) Service {
	return &service{
		sectionRepository:            sectionRepo,
		dicesRepository:              dicesRepo,
		battleRepository:             battleRepo,
		playerRepository:             playerRepo,
		playerSectionEnemyRepository: playerSectionEnemyRepo,
		db:                           db,
		logService:                   logService,
	}
}

func (s *service) Action(ctx context.Context, player entities.Player, req gameDTO.ActionRequest) (gameDTO.ActionResponse, error) {
	var result gameDTO.ActionResult

	battleResult, err := s.Battle(ctx, player, req)
	if err != nil {
		return dto.ActionResponse{}, err
	}

	if battleResult.Finish {
		result = gameDTO.ResultTrue
	}

	return dto.ActionResponse{
		Result: result,
	}, nil
}

func (s *service) Battle(ctx context.Context, player entities.Player, req dto.ActionRequest) (dto.BattleDto, error) {
	var logService log.PlayerLogService
	if s.logService != nil {
		logService = s.logService
	}
	common, err := battleCommon.NewCommonWithRepositories(ctx, s.db, &player, s.battleRepository, s.playerSectionEnemyRepository, logService)
	if err != nil {
		return dto.BattleDto{}, err
	}

	battle, err := common.Action(&req.Weapon)
	if err != nil {
		return dto.BattleDto{}, err
	}

	_, err = s.battleRepository.AddRecord(s.db, battle)
	if err != nil {
		return dto.BattleDto{}, err
	}

	return dto.BattleDto{
		Finish: true,
	}, nil
}
