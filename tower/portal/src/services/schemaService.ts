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
   * Fetches a schema from the server or cache
   * @param schemaType - The type of schema to fetch (e.g., 'woocommerce', 'brincr')
   * @returns The schema data
   */
  async getSchema(schemaType: string): Promise<Record<string, any>> {
    // Check cache first
    if (this.cache[schemaType]) {
      return this.cache[schemaType];
    }

    try {
      // In a real app, this would be a fetch to your server
      // For now, we'll use a simulated fetch
      const response = await this.simulateFetch(schemaType);
      
      // Cache the schema
      this.cache[schemaType] = response;
      
      return response;
    } catch (error) {
      console.error(`Error fetching ${schemaType} schema:`, error);
      throw error;
    }
  }

    /**
     * Fetches source/target fields and suggested mappings from the backend
     */
    async getConnectionSchemas(sourceType: string, targetType: string): Promise<{
        sourceFields: any[];
        targetFields: any[];
        mappings: MappingDefinition[];
    }> {
        try {
            const response = await fetch(`http://localhost:8080/api/mappings/schema?source=${sourceType}&target=${targetType}`);

            if (!response.ok) {
                throw new Error(`Failed to fetch backend schema mapping: ${response.statusText}`);
            }

            const data = await response.json();

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
   * Simulates a fetch to the server
   * In a real app, this would be replaced with actual fetch calls
   * @param schemaType - The type of schema to fetch
   * @returns The schema data
   */
  private async simulateFetch(schemaType: string): Promise<Record<string, any>> {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 500));
    
    // Return mock schemas
    if (schemaType === 'woocommerce') {
      return {
        "id": 794,
        "name": "Premium Quality",
        "slug": "premium-quality-19",
        "permalink": "https://example.com/product/premium-quality-19/",
        "date_created": "2017-03-23T17:01:14",
        "date_modified": "2017-03-23T17:01:14",
        "type": "simple",
        "status": "publish",
        "featured": false,
        "catalog_visibility": "visible",
        "description": "<p>Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.</p>",
        "short_description": "<p>Pellentesque habitant morbi tristique senectus et netus.</p>",
        "sku": "PROD-123",
        "price": "21.99",
        "regular_price": "21.99",
        "sale_price": "",
        "on_sale": false,
        "purchasable": true,
        "total_sales": 0,
        "virtual": false,
        "downloadable": false,
        "manage_stock": false,
        "stock_quantity": null,
        "stock_status": "instock",
        "weight": "",
        "dimensions": {
          "length": "10",
          "width": "20",
          "height": "15"
        },
        "shipping_required": true,
        "shipping_taxable": true,
        "shipping_class": "",
        "reviews_allowed": true,
        "average_rating": "0.00",
        "rating_count": 0,
        "categories": [
          {
            "id": 9,
            "name": "Clothing",
            "slug": "clothing"
          },
          {
            "id": 14,
            "name": "T-shirts",
            "slug": "t-shirts"
          }
        ],
        "tags": [],
        "images": [
          {
            "id": 792,
            "src": "https://example.com/wp-content/uploads/2017/03/T_2_front-4.jpg",
            "name": "",
            "alt": ""
          }
        ],
        "attributes": [],
        "variations": [],
        "meta_data": []
      };
    } else if (schemaType === 'brincr') {
      return {
                "product_type": "default",
                "category_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "vat_rate_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "unit_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "matrix_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                "user_defined_1": "string",
                "user_defined_2": "string",
                "user_defined_3": "string",
                "user_defined_4": "string",
                "user_defined_5": "string",
                "user_defined_6": "string",
                "user_defined_7": "string",
                "user_defined_8": "string",
                "user_defined_9": "string",
                "user_defined_10": "string",
                "stock": {
                    "settings": {
                        "is_stock_item": true,
                        "critical_level": 9999999.9999,
                        "default_sales_warehouse_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                        "default_purchase_warehouse_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
                    }
                },
                "parts": [
                    {
                        "product_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                        "amount": 9999999.9999
                    }
                ],
                "assemble_type": "assemble_on_order",
                "code": "string",
                "description": {
                    "en_GB": "text",
                    "nl_NL": "dutch text",
                    "fr_FR": "french text"
                },
                "description_extended": {
                    "en_GB": "text",
                    "nl_NL": "dutch text",
                    "fr_FR": "french text"
                },
                "price_excl_vat": 9999999.9999,
                "image": "string",
                "file": "string",
                "note": "string",
                "stock_level": 9999999.9999,
                "barcode": "string",
                "package_length": 999999.99,
                "package_width": 999999.99,
                "package_height": 999999.99,
                "package_volume": 999999.99,
                "weight": 99999.999,
                "fixed_stock_price": 99999.999,
                "manufacturers_suggested_retail_price": 99999.999,
                "manufacturers_suggested_retail_price_type": "Single MSRP price",
                "accountancy": {
                    "export_override_income": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "export_override_income_eu": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                    "export_override_income_non_eu": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
                },
                "web_portal": {
                    "description": {
                        "en_GB": "english text",
                        "nl_NL": "dutch text",
                        "de_DE": "german text"
                    }
                },
                "is_alert_message_shown": true,
                "alert_message": "string",
                "web_shop_uri": "string",
                "is_from_web_shop": true,
                "is_active": true,
                "is_shown_in_price_list": true,
                "is_shown_on_labels": true,
                "is_shown_on_packing_slip": true,
                "is_purchase_price_customizable": true
      };
    }
    
    throw new Error(`Unknown schema type: ${schemaType}`);
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
