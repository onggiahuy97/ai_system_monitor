package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var prevNetStats map[string]net.IOCountersStat

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func getGPUUsage() string {
	// For macOS, we'll use ioreg to get GPU info
	// This is a simplified version - GPU monitoring on macOS is limited without special tools
	cmd := exec.Command("ioreg", "-r", "-d", "1", "-w", "0", "-c", "IOAccelerator")
	output, err := cmd.Output()
	if err != nil {
		return "GPU: N/A (requires Metal-compatible GPU)"
	}

	// Basic detection - more detailed GPU stats require Metal API or powermetrics
	if strings.Contains(string(output), "IOAccelerator") {
		return "GPU: Active (detailed stats require powermetrics with sudo)"
	}
	return "GPU: N/A"
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getNetworkStats() (string, error) {
	counters, err := net.IOCounters(true)
	if err != nil {
		return "", err
	}

	currentStats := make(map[string]net.IOCountersStat)
	for _, counter := range counters {
		currentStats[counter.Name] = counter
	}

	var result strings.Builder
	result.WriteString("Network Usage:\n")

	if prevNetStats == nil {
		prevNetStats = currentStats
		result.WriteString("  (Calculating...)")
		return result.String(), nil
	}

	var totalSent, totalRecv uint64
	for name, current := range currentStats {
		if prev, ok := prevNetStats[name]; ok {
			sent := current.BytesSent - prev.BytesSent
			recv := current.BytesRecv - prev.BytesRecv

			// Only show active interfaces
			if sent > 0 || recv > 0 {
				totalSent += sent
				totalRecv += recv
			}
		}
	}

	// Convert to per-second rates (5 second interval)
	sentPerSec := totalSent / 5
	recvPerSec := totalRecv / 5

	result.WriteString(fmt.Sprintf("  ↑ Upload:   %s/s\n", formatBytes(sentPerSec)))
	result.WriteString(fmt.Sprintf("  ↓ Download: %s/s", formatBytes(recvPerSec)))

	prevNetStats = currentStats
	return result.String(), nil
}

func displayStats() {
	clearScreen()

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           MacBook System Monitor - Live Stats               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// CPU Usage
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Printf("CPU: Error - %v\n", err)
	} else {
		cpuUsage := percentages[0]
		cpuCount, _ := cpu.Counts(true)
		fmt.Printf("🖥️  CPU Usage: %.2f%% (%d cores)\n", cpuUsage, cpuCount)

		// Visual bar
		barLength := int(cpuUsage / 2)
		bar := strings.Repeat("█", barLength) + strings.Repeat("░", 50-barLength)
		fmt.Printf("   [%s]\n", bar)
	}

	fmt.Println()

	// Memory Usage
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("RAM: Error - %v\n", err)
	} else {
		fmt.Printf("💾 RAM Usage: %.2f%% (Used: %s / Total: %s)\n",
			vmStat.UsedPercent,
			formatBytes(vmStat.Used),
			formatBytes(vmStat.Total))

		// Visual bar
		barLength := int(vmStat.UsedPercent / 2)
		bar := strings.Repeat("█", barLength) + strings.Repeat("░", 50-barLength)
		fmt.Printf("   [%s]\n", bar)
	}

	fmt.Println()

	// GPU Usage
	gpuInfo := getGPUUsage()
	fmt.Printf("🎮 %s\n", gpuInfo)

	fmt.Println()

	// Network Stats
	netStats, err := getNetworkStats()
	if err != nil {
		fmt.Printf("📡 Network: Error - %v\n", err)
	} else {
		fmt.Printf("📡 %s\n", netStats)
	}

	fmt.Println()
	fmt.Printf("⏱️  Last updated: %s\n", time.Now().Format("15:04:05"))
	fmt.Println("\n[Press Ctrl+C to exit]")
}

func main() {
	if runtime.GOOS != "darwin" {
		fmt.Println("This program is designed for macOS (darwin)")
		os.Exit(1)
	}

	fmt.Println("Starting system monitor...")
	time.Sleep(time.Second)

	// Initial display
	displayStats()

	// Update every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		displayStats()
	}
}
