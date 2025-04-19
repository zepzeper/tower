import { MappingDefinition } from '../utils/schemaTransformer';

export interface ConnectionSchemas {
    sourceSchema: Record<string, any>;
    targetSchema: Record<string, any>;
    sourceType: string;
    targetType: string;
}

export interface SavedMapping {
    connectionId: string;
    sourceType: string;
    targetType: string;
    mappings: MappingDefinition[];
    mappingData: MappingMetadata[];
}

export interface MappingMetadata {
    id: string;
    sourcePath: string | undefined;
    targetPath: string | undefined;
    transform: string | null;
    sourceFieldId: string;
    targetFieldId: string;
}

/**
 * Service for fetching and caching API schemas
 */
class SchemaService {
    private cache: Record<string, Record<string, any>> = {};

    /**
     * Fetches source/target fields and suggested mappings from the backend
     */
    async getConnectionSchemas(sourceType: string, targetType: string): Promise<{
        sourceFields: any[];
        targetFields: any[];
        mappings: MappingDefinition[];
    }> {
        try {
            const response = await fetch(`/api/mappings/schema?source=${sourceType}&target=${targetType}&operation=products`);

            if (!response.ok) {
                throw new Error(`Failed to fetch backend schema mapping: ${response.statusText}`);
            }

            const json = await response.json();
            const data = json.data;

            return {
                sourceFields: data.sourceFields,
                targetFields: data.targetFields,
                mappings: data.mappings,
            };
        } catch (error) {
            console.error('Error fetching backend schema mapping:', error);
            throw error;
        }
    }

    /**
     * Tests a mapping configuration with real data
     * @param connectionId - The ID of the connection
     * @param sourceType - Source system type
     * @param targetType - Target system type
     * @param mappings - The mappings to test
     * @returns Source data and transformed result
     */
    async testMapping(
        connectionId: string, 
        sourceType: string, 
        targetType: string, 
        mappings: MappingMetadata[]
    ): Promise<{sourceData: any, transformedData: any}> {
        try {
            const response = await fetch(`/api/mappings/test`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    connectionId,
                    sourceType,
                    targetType,
                    mappings
                })
            });

            if (!response.ok) {
                throw new Error(`Failed to test mappings: ${response.statusText}`);
            }

            return await response.json();
        } catch (error) {
            console.error('Error testing mappings:', error);
            throw error;
        }
    }

    /**
   * Fetches any saved mapping configuration for this connection
   * @param connectionId - The ID of the connection
   * @returns The saved mapping or null if none exists
   */
    async getSavedMapping(connectionId: string): Promise<SavedMapping | null> {
        try {
            // In a real app, this would fetch from your server
            // Simulate a response for now
            await new Promise(resolve => setTimeout(resolve, 300));

            // For demo purposes, return null to indicate no saved mapping
            // In a real app, you'd return actual saved mappings if they exist
            return null;
        } catch (error) {
            console.error('Error fetching saved mapping:', error);
            return null;
        }
    }

    /**
   * Saves a mapping configuration
   * @param connectionId - The ID of the connection
   * @param mapping - The mapping configuration to save
   * @returns Success indicator
   */
    async saveMapping(connectionId: string, mapping: SavedMapping): Promise<boolean> {
        try {
            // In a real app, this would save to your server
            console.log('Saving mapping for connection:', connectionId, mapping);

            // Simulate successful save
            await new Promise(resolve => setTimeout(resolve, 500));
            return true;
        } catch (error) {
            console.error('Error saving mapping:', error);
            throw error;
        }
    }
}

export default new SchemaService();
