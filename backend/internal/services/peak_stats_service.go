package services

import (
        "context"
        "fmt"
        "log"
        "time"

        "historical-events-backend/internal/database/repositories"
)

// PeakStatsService tracks peak concurrent sessions in the background
type PeakStatsService struct {
        userRepo *repositories.UserRepository
        logger   *log.Logger
}

// NewPeakStatsService creates a new PeakStatsService
func NewPeakStatsService(userRepo *repositories.UserRepository, logger *log.Logger) *PeakStatsService {
        return &PeakStatsService{
                userRepo: userRepo,
                logger:   logger,
        }
}

// Start begins the peak stats tracking loop
func (s *PeakStatsService) Start(ctx context.Context) {
        ticker := time.NewTicker(60 * time.Second)
        defer ticker.Stop()

        s.logger.Println("Peak stats service started")

        for {
                select {
                case <-ticker.C:
                        if err := s.updatePeakStats(); err != nil {
                                s.logger.Printf("Error updating peak stats: %v", err)
                        }
                case <-ctx.Done():
                        s.logger.Println("Peak stats service stopped")
                        return
                }
        }
}

// updatePeakStats counts active sessions and updates peak if higher
func (s *PeakStatsService) updatePeakStats() error {
        currentCount, err := s.userRepo.CountAllActiveSessions()
        if err != nil {
                return fmt.Errorf("failed to count active sessions: %w", err)
        }

        currentPeak, err := s.userRepo.GetPeakConcurrentSessions()
        if err != nil {
                return fmt.Errorf("failed to get current peak: %w", err)
        }

        if currentCount > currentPeak {
                err = s.userRepo.UpdatePeakIfHigher(currentCount)
                if err != nil {
                        return fmt.Errorf("failed to update peak: %w", err)
                }
                s.logger.Printf("New peak concurrent sessions: %d (previous: %d)", currentCount, currentPeak)
        }

        return nil
}
