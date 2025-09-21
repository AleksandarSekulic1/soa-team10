import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MonitoringService, MonitoringDashboard, ExternalTools } from '../../services/monitoring.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-monitoring',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './monitoring.component.html',
  styleUrls: ['./monitoring.component.scss']
})
export class MonitoringComponent implements OnInit, OnDestroy {
  dashboard$: Observable<MonitoringDashboard | null>;
  isLoading$: Observable<boolean>;
  error$: Observable<string | null>;
  
  externalTools: ExternalTools;

  constructor(private monitoringService: MonitoringService) {
    this.dashboard$ = this.monitoringService.dashboard$;
    this.isLoading$ = this.monitoringService.isLoading$;
    this.error$ = this.monitoringService.error$;
    this.externalTools = this.monitoringService.getExternalTools();
  }

  ngOnInit(): void {
    // Start monitoring when component initializes
    console.log("🚀 Monitoring component initializing...");
    this.monitoringService.debugDashboardState();
    this.monitoringService.startMonitoring();
    
    // Debug subscription to see when dashboard data changes
    this.dashboard$.subscribe(data => {
      console.log("📊 Dashboard data received in component:", data);
    });
  }

  ngOnDestroy(): void {
    // Stop monitoring when component is destroyed
    this.monitoringService.stopMonitoring();
  }

  /**
   * Manual refresh of monitoring data
   */
  onRefresh(): void {
    this.monitoringService.refresh().subscribe();
  }

  /**
   * Open external monitoring tool in new tab
   */
  openExternalTool(url: string): void {
    window.open(url, '_blank');
  }

  /**
   * Get status class for styling
   */
  getStatusClass(status: string | undefined): string {
    return status || 'unknown';
  }

  /**
   * Type guard to check if dashboard is valid
   */
  isDashboard(dashboard: any): dashboard is MonitoringDashboard {
    return dashboard && typeof dashboard === 'object' && 'status' in dashboard;
  }

  /**
   * Check if system is healthy
   */
  isSystemHealthy(): boolean {
    return this.monitoringService.isSystemHealthy();
  }
}