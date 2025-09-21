import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, interval, tap, catchError, of } from 'rxjs';
import { environment } from '../../environments/environment';

export interface GatewayMetrics {
  service: string;
  uptime: string;
  total_requests: number;
  total_errors: number;
  active_requests: number;
  average_response_time_ms: number;
  memory_usage_mb: string;
  goroutines: number;
  timestamp: string;
}

export interface ContainerMetrics {
  cpu_usage_percent: number;
  memory_usage_percent: number;
  disk_usage_percent: number;
  network_rx_bytes: number;
  network_tx_bytes: number;
}

export interface HostMetrics {
  cpu_usage_percent: number;
  memory_usage_percent: number;
  disk_usage_percent: number;
  network_rx_bytes: number;
  network_tx_bytes: number;
}

export interface SystemMetrics {
  container: ContainerMetrics;
  host: HostMetrics;
}

export interface MonitoringDashboard {
  status: 'healthy' | 'unhealthy' | 'unknown';
  gateway: GatewayMetrics;
  system: SystemMetrics;
}

export interface ExternalTools {
  prometheus: string;
  grafana: string;
  nodeExporter: string;
  cadvisor: string;
}

@Injectable({
  providedIn: 'root'
})
export class MonitoringService {
  private apiUrl = `${environment.apiUrl}`;

  private dashboardSubject = new BehaviorSubject<MonitoringDashboard | null>(null);
  dashboard$ = this.dashboardSubject.asObservable();

  private isLoadingSubject = new BehaviorSubject<boolean>(false);
  isLoading$ = this.isLoadingSubject.asObservable();

  private errorSubject = new BehaviorSubject<string | null>(null);
  error$ = this.errorSubject.asObservable();

  private refreshInterval = 4000; // 4 seconds
  private intervalSubscription: any;

  constructor(private http: HttpClient) {}

  /**
   * Start automatic monitoring data refresh
   */
  startMonitoring(): void {
    console.log('🚀 Starting monitoring service...');
    this.loadDashboard().subscribe({
      next: (data) => console.log('🔄 Initial monitoring data loaded:', data),
      error: (error) => console.error('❌ Initial monitoring load failed:', error)
    });
    
    if (this.intervalSubscription) {
      this.intervalSubscription.unsubscribe();
    }
    
    this.intervalSubscription = interval(this.refreshInterval).subscribe(() => {
      console.log('⏰ Refreshing monitoring data...');
      this.loadDashboard().subscribe({
        next: (data) => console.log('🔄 Monitoring data refreshed:', data),
        error: (error) => console.error('❌ Monitoring refresh failed:', error)
      });
    });
  }

  /**
   * Stop automatic monitoring data refresh
   */
  stopMonitoring(): void {
    if (this.intervalSubscription) {
      this.intervalSubscription.unsubscribe();
      this.intervalSubscription = null;
    }
  }

  /**
   * Load monitoring dashboard data from API Gateway
   */
  loadDashboard(): Observable<MonitoringDashboard> {
    this.isLoadingSubject.next(true);
    this.errorSubject.next(null);
    
    console.log('🔍 Attempting to load monitoring dashboard from:', `${this.apiUrl}/monitoring`);
    return this.http.get<MonitoringDashboard>(`${this.apiUrl}/monitoring`).pipe(
      tap(data => {
        console.log('✅ Monitoring dashboard data loaded successfully:', data);
        console.log('📊 Assigning data to dashboard subject...');
        this.dashboardSubject.next(data);
        this.isLoadingSubject.next(false);
        console.log('🎯 Dashboard subject current value:', this.dashboardSubject.value);
      }),
      catchError(error => {
        console.error('❌ Failed to load monitoring data:', error);
        console.error('❌ Error details:', {
          status: error.status,
          statusText: error.statusText,
          url: error.url,
          message: error.message
        });
        this.errorSubject.next('Failed to load monitoring data. Please check if the API Gateway is running.');
        this.isLoadingSubject.next(false);
        
        // Create a fallback dashboard
        const fallbackDashboard: MonitoringDashboard = {
          status: 'unknown',
          gateway: {
            service: 'api-gateway',
            uptime: 'N/A',
            total_requests: 0,
            total_errors: 0,
            active_requests: 0,
            average_response_time_ms: 0,
            memory_usage_mb: 'N/A',
            goroutines: 0,
            timestamp: new Date().toISOString()
          },
          system: {
            container: {
              cpu_usage_percent: 0,
              memory_usage_percent: 0,
              disk_usage_percent: 0,
              network_rx_bytes: 0,
              network_tx_bytes: 0
            },
            host: {
              cpu_usage_percent: 0,
              memory_usage_percent: 0,
              disk_usage_percent: 0,
              network_rx_bytes: 0,
              network_tx_bytes: 0
            }
          }
        };
        
        // Only set fallback data if there's no existing data
        if (!this.dashboardSubject.value) {
          console.log('🆘 Setting fallback dashboard data');
          this.dashboardSubject.next(fallbackDashboard);
        }
        
        return of(fallbackDashboard);
      })
    );
  }

  /**
   * Get current dashboard data
   */
  getCurrentDashboard(): MonitoringDashboard | null {
    const current = this.dashboardSubject.value;
    console.log('🔍 getCurrentDashboard called, returning:', current);
    return current;
  }

  /**
   * Debug method to check dashboard state
   */
  debugDashboardState(): void {
    console.log('🐛 DEBUG Dashboard State:');
    console.log('  - Dashboard Subject Value:', this.dashboardSubject.value);
    console.log('  - Is Loading:', this.isLoadingSubject.value);
    console.log('  - Error:', this.errorSubject.value);
    console.log('  - API URL:', this.apiUrl);
  }

  /**
   * Check if monitoring system is healthy
   */
  isSystemHealthy(): boolean {
    const dashboard = this.getCurrentDashboard();
    return dashboard?.status === 'healthy';
  }

  /**
   * Get external monitoring tool URLs
   */
  getExternalTools(): ExternalTools {
    return {
      prometheus: 'http://localhost:9090',
      grafana: 'http://localhost:3000',
      nodeExporter: 'http://localhost:9100/metrics',
      cadvisor: 'http://localhost:8087' // Updated port after your change
    };
  }

  /**
   * Set refresh interval for monitoring data
   */
  setRefreshInterval(intervalMs: number): void {
    this.refreshInterval = intervalMs;
    
    // Restart monitoring with new interval if it's currently running
    if (this.intervalSubscription) {
      this.stopMonitoring();
      this.startMonitoring();
    }
  }

  /**
   * Force refresh monitoring data
   */
  refresh(): Observable<MonitoringDashboard> {
    return this.loadDashboard();
  }

  /**
   * Clean up when service is destroyed
   */
  ngOnDestroy(): void {
    this.stopMonitoring();
  }
}