class CacheManager {
  private cache = new Map<string, { value: any; timestamp: number }>();
  private dependencies = new Map<string, Set<string>>();
  private defaultTTL = 5 * 60 * 1000; // 5 minutes default

  /**
   * Store a value in the cache
   * @param key - Unique cache key
   * @param value - The value to cache
   * @param deps - Array of resource IDs this cache entry depends on
   * @param ttl - Time to live in milliseconds (optional)
   */
  set(key: string, value: any, deps: string[] = [], ttl: number = this.defaultTTL) {
    this.cache.set(key, {
      value,
      timestamp: Date.now() + ttl
    });

    // Register dependencies
    deps.forEach(dep => {
      if (!this.dependencies.has(dep)) {
        this.dependencies.set(dep, new Set());
      }
      this.dependencies.get(dep)?.add(key);
    });
  }

  /**
   * Get a value from the cache
   * @param key - Cache key
   * @returns The cached value or undefined if not found or expired
   */
  get(key: string) {
    const cached = this.cache.get(key);

    if (!cached) return undefined;

    // Check if entry is expired
    if (cached.timestamp < Date.now()) {
      this.cache.delete(key);
      return undefined;
    }

    return cached.value;
  }

  /**
   * Check if a key exists in the cache and is not expired
   * @param key - Cache key
   * @returns Boolean indicating if the key exists and is valid
   */
  has(key: string): boolean {
    const cached = this.cache.get(key);
    if (!cached) return false;

    if (cached.timestamp < Date.now()) {
      this.cache.delete(key);
      return false;
    }

    return true;
  }

  /**
   * Invalidate all caches dependent on a resource
   * @param resourceId - ID of the resource that changed
   */
  invalidateResource(resourceId: string) {
    const dependentKeys = this.dependencies.get(resourceId);
    if (dependentKeys) {
      dependentKeys.forEach(key => {
        this.cache.delete(key);
      });
      // Clear the dependency tracking for this resource
      this.dependencies.delete(resourceId);
    }
  }

  /**
   * Invalidate multiple resources at once
   * @param resourceIds - Array of resource IDs to invalidate
   */
  invalidateResources(resourceIds: string[]) {
    resourceIds.forEach(id => this.invalidateResource(id));
  }

  /**
   * Invalidate a specific cache key
   * @param key - Cache key to invalidate
   */
  invalidate(key: string) {
    this.cache.delete(key);
  }

  /**
   * Clear all cache entries
   */
  clear() {
    this.cache.clear();
    this.dependencies.clear();
  }

  /**
   * Get all dependency relationships for debugging
   */
  getDependencyMap() {
    const map: Record<string, string[]> = {};
    this.dependencies.forEach((keys, resource) => {
      map[resource] = Array.from(keys);
    });
    return map;
  }
}

export default CacheManager;
