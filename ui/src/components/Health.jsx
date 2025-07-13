import React, { useState, useEffect } from 'react'
import { adminAPI } from '../api/admin'
import { 
  Activity, 
  Database, 
  Server, 
  Clock, 
  CheckCircle, 
  AlertCircle,
  RefreshCw,
  Trash2,
  HardDrive,
  Monitor
} from 'lucide-react';

const Health = () => {
  const [health, setHealth] = useState(null);
  const [systemInfo, setSystemInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [purging, setPurging] = useState(false);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [healthData, systemData] = await Promise.all([
        adminAPI.getHealth(),
        adminAPI.getSystemInfo()
      ]);
      setHealth(healthData);
      setSystemInfo(systemData);
    } catch (error) {
      console.error('Failed to fetch health data:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    
    // Auto-refresh every 30 seconds
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, []);

  const handlePurgeLogs = async (days) => {
    if (window.confirm(`Are you sure you want to purge logs older than ${days} days?`)) {
      try {
        setPurging(true);
        const result = await adminAPI.purgeOldLogs(days);
        alert(`Successfully purged ${result.deleted_count} log entries`);
        fetchData(); // Refresh data
      } catch (error) {
        alert('Failed to purge logs: ' + error.message);
      } finally {
        setPurging(false);
      }
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-2xl font-bold text-gray-900">System Health</h1>
          <p className="mt-1 text-sm text-gray-500">
            Monitor system health and manage resources
          </p>
        </div>
        <div className="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
          <button
            onClick={fetchData}
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700"
          >
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh
          </button>
        </div>
      </div>

      {/* Overall Health Status */}
      {health && (
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                {health.status === 'healthy' ? (
                  <CheckCircle className="h-12 w-12 text-green-400" />
                ) : (
                  <AlertCircle className="h-12 w-12 text-red-400" />
                )}
              </div>
              <div className="ml-5">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                  System Status: {health.status === 'healthy' ? 'Healthy' : 'Unhealthy'}
                </h3>
                <div className="mt-2 grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm text-gray-500">
                  <div className="flex items-center space-x-2">
                    <Clock className="h-4 w-4" />
                    <span>Uptime: {health.uptime}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Monitor className="h-4 w-4" />
                    <span>Version: {health.version}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Activity className="h-4 w-4" />
                    <span>Last Check: {new Date(health.timestamp).toLocaleTimeString()}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Health Checks */}
      {health?.checks && (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <div className="px-4 py-5 sm:px-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              Health Checks
            </h3>
            <p className="mt-1 max-w-2xl text-sm text-gray-500">
              Individual component health status
            </p>
          </div>
          <ul className="divide-y divide-gray-200">
            {Object.entries(health.checks).map(([name, check]) => (
              <li key={name} className="px-4 py-4 sm:px-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      {check.status === 'healthy' ? (
                        <CheckCircle className="h-6 w-6 text-green-400" />
                      ) : (
                        <AlertCircle className="h-6 w-6 text-red-400" />
                      )}
                    </div>
                    <div className="ml-4">
                      <div className="flex items-center space-x-2">
                        <p className="text-sm font-medium text-gray-900 capitalize">
                          {name}
                        </p>
                        <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                          check.status === 'healthy' 
                            ? 'bg-green-100 text-green-800' 
                            : 'bg-red-100 text-red-800'
                        }`}>
                          {check.status}
                        </span>
                      </div>
                      {check.message && (
                        <p className="text-sm text-gray-500">{check.message}</p>
                      )}
                    </div>
                  </div>
                  {check.latency && (
                    <div className="text-sm text-gray-500">
                      {check.latency.toFixed(2)}ms
                    </div>
                  )}
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Database Information */}
      {systemInfo?.database && (
        <div className="bg-white overflow-hidden shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <div className="flex items-center mb-4">
              <Database className="h-6 w-6 text-blue-500 mr-2" />
              <h3 className="text-lg leading-6 font-medium text-gray-900">
                Database Information
              </h3>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <div className="space-y-3">
                  <div>
                    <label className="text-sm font-medium text-gray-500">Version</label>
                    <p className="text-sm text-gray-900">{systemInfo.database.version}</p>
                  </div>
                  <div>
                    <label className="text-sm font-medium text-gray-500">Database Size</label>
                    <p className="text-sm text-gray-900">{systemInfo.database.size}</p>
                  </div>
                </div>
              </div>
              
              {systemInfo.database.table_sizes && (
                <div>
                  <label className="text-sm font-medium text-gray-500 mb-2 block">Table Sizes</label>
                  <div className="space-y-2">
                    {systemInfo.database.table_sizes.map((table, index) => (
                      <div key={index} className="flex justify-between items-center">
                        <span className="text-sm text-gray-900">{table.table_name}</span>
                        <span className="text-sm text-gray-500">{table.size}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Log Management */}
      <div className="bg-white overflow-hidden shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center mb-4">
            <HardDrive className="h-6 w-6 text-yellow-500 mr-2" />
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              Log Management
            </h3>
          </div>
          
          <p className="text-sm text-gray-500 mb-4">
            Purge old log entries to free up database space. This action cannot be undone.
          </p>
          
          <div className="flex flex-wrap gap-3">
            {[7, 30, 90, 365].map(days => (
              <button
                key={days}
                onClick={() => handlePurgeLogs(days)}
                disabled={purging}
                className="inline-flex items-center px-3 py-2 border border-red-300 text-sm leading-4 font-medium rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
              >
                <Trash2 className="h-4 w-4 mr-1" />
                Purge {days}d+
              </button>
            ))}
          </div>
          
          {purging && (
            <div className="mt-4 flex items-center space-x-2 text-sm text-gray-600">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
              <span>Purging logs...</span>
            </div>
          )}
        </div>
      </div>

      {/* System Resources */}
      <div className="bg-white overflow-hidden shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center mb-4">
            <Server className="h-6 w-6 text-green-500 mr-2" />
            <h3 className="text-lg leading-6 font-medium text-gray-900">
              System Resources
            </h3>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-500">Memory</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {health?.checks?.memory?.status === 'healthy' ? 'Normal' : 'Warning'}
                  </p>
                </div>
                <Activity className="h-8 w-8 text-blue-500" />
              </div>
            </div>
            
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-500">Database</p>
                  <p className="text-2xl font-bold text-gray-900">
                    {health?.checks?.database?.status === 'healthy' ? 'Connected' : 'Error'}
                  </p>
                </div>
                <Database className="h-8 w-8 text-green-500" />
              </div>
              {health?.checks?.database?.latency && (
                <p className="text-xs text-gray-500 mt-1">
                  Latency: {health.checks.database.latency.toFixed(2)}ms
                </p>
              )}
            </div>
            
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-500">Uptime</p>
                  <p className="text-lg font-bold text-gray-900">
                    {health?.uptime || 'Unknown'}
                  </p>
                </div>
                <Clock className="h-8 w-8 text-purple-500" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Health