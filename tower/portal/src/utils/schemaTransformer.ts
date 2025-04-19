export interface FieldDefinition {
  id: string;
  name: string;
  type: string;
  path: string;
  sample?: string;
  required?: boolean;
}

export interface MappingDefinition {
  id: string;
  sourceField: string;
  targetField: string;
  transform: string | null;
}

export interface MappingData {
  sourceFields: FieldDefinition[];
  targetFields: FieldDefinition[];
  mappings: MappingDefinition[];
}

export interface SchemaTransformOptions {
  prefix?: string;
  includeRequired?: boolean;
}

/**
 * Transforms the raw API schema into a format suitable for the field mapping component
 * @param schema - The API schema object 
 * @param type - Either 'source' or 'target' to identify the schema type
 * @param options - Additional options for transformation
 * @returns Formatted fields for the mapping component
 */
export const transformSchema = (
  schema: Record<string, any>, 
  type: 'source' | 'target', 
  options: SchemaTransformOptions = {}
): FieldDefinition[] => {
  const { prefix = type === 'source' ? 's' : 't', includeRequired = true } = options;
  const fields: FieldDefinition[] = [];
  let counter = 1;

  // Function to recursively process schema objects
  const processObject = (obj: any, path = '', parentName = ''): void => {
    if (!obj || typeof obj !== 'object') return;

    // Handle arrays in the schema (for nested objects)
    if (Array.isArray(obj) && obj.length > 0) {
      // If it's an array of objects, process the first item as an example
      if (typeof obj[0] === 'object') {
        processObject(obj[0], `${path}[]`, parentName);
      }
      return;
    }

    // Process each property in the object
    Object.entries(obj).forEach(([key, value]) => {
      const currentPath = path ? `${path}.${key}` : key;
      const displayName = parentName ? `${parentName}.${key}` : key;
      
      if (value === null || value === undefined) {
        // Skip null values
        return;
      }
      
      // Determine the field type
      let fieldType = 'string';
      let sampleValue = '';
      
      if (typeof value === 'number') {
        fieldType = 'number';
        sampleValue = String(value);
      } else if (typeof value === 'boolean') {
        fieldType = 'boolean';
        sampleValue = String(value);
      } else if (typeof value === 'string') {
        fieldType = 'string';
        sampleValue = value || `Sample ${key}`;
      } else if (Array.isArray(value)) {
        fieldType = 'array';
        
        // For arrays, determine the type of items
        if (value.length > 0) {
          const firstItem = value[0];
          if (typeof firstItem === 'object' && firstItem !== null) {
            // For arrays of objects
            fieldType = 'array.object';
            
            // Process array item properties
            Object.keys(firstItem).forEach(itemKey => {
              processObject(
                { [itemKey]: firstItem[itemKey] }, 
                `${currentPath}[].${itemKey}`, 
                `${displayName}[].${itemKey}`
              );
            });
            
            // Create a sample showing array format
            if (Object.keys(firstItem).length > 0) {
              const sampleKey = Object.keys(firstItem)[0];
              sampleValue = `[{"${sampleKey}": "${firstItem[sampleKey]}"}]`;
            } else {
              sampleValue = '[]';
            }
          } else {
            // For arrays of primitives
            fieldType = `array.${typeof firstItem}`;
            sampleValue = JSON.stringify([firstItem]);
          }
        } else {
          sampleValue = '[]';
        }
      } else if (typeof value === 'object') {
        // For objects, process nested properties
        fieldType = 'object';
        processObject(value, currentPath, displayName);
        
        // If it's a simple object with few properties, create a sample
        const keys = Object.keys(value);
        if (keys.length > 0 && keys.length <= 3) {
          sampleValue = JSON.stringify(value);
          if (sampleValue.length > 50) {
            sampleValue = sampleValue.substring(0, 47) + '...';
          }
        } else {
          sampleValue = '{...}';
        }
      }
      
      // Add the field to our fields array if it's a primitive or array
      if (fieldType !== 'object' || Object.keys(value).length === 0) {
        // Determine if field is required (simplistic approach for demo)
        const isRequired = includeRequired && 
          (key === 'id' || key === 'sku' || key === 'name' || key === 'price' || 
           key === 'code' || key === 'description');
           
        fields.push({
          id: `${prefix}${counter++}`,
          name: key,
          type: fieldType,
          path: currentPath.replace(/\.\[\]\./g, '[].'),
          sample: sampleValue.substring(0, 100), // Limit sample size
          ...(type === 'target' && includeRequired ? { required: isRequired } : {})
        });
      }
    });
  };

  processObject(schema);
  return fields;
};

/**
 * Transforms API schemas and generates predefined mappings
 * @param sourceSchema - The source API schema
 * @param targetSchema - The target API schema
 * @returns Object containing sourceFields, targetFields, and mappings
 */
export const generateMappingData = (
  sourceSchema: Record<string, any>, 
  targetSchema: Record<string, any>
): MappingData => {
  // Transform schemas to field arrays
  const sourceFields = transformSchema(sourceSchema, 'source');
  const targetFields = transformSchema(targetSchema, 'target');
  
  // Generate automatic mappings based on field similarities
  const mappings: MappingDefinition[] = [];
  let mappingCounter = 1;
  
  // Find potential field matches
  targetFields.forEach(targetField => {
    const targetName = targetField.name.toLowerCase();
    
    // Look for matching source field by name
    const sourceField = sourceFields.find(source => {
      const sourceName = source.name.toLowerCase();
      
      // Direct matches
      if (sourceName === targetName) return true;
      
      // Check for partial matches
      if (targetName.includes(sourceName) || sourceName.includes(targetName)) return true;
      
      // Check for common mappings by convention
      if (
        (targetName === 'sku' && sourceName === 'id') ||
        (targetName === 'description' && sourceName === 'short_description') ||
        (targetName === 'price' && sourceName === 'regular_price') ||
        (targetName === 'image' && sourceName === 'images')
      ) {
        return true;
      }
      
      return false;
    });
    
    if (sourceField) {
      // Determine if transformation is needed
      let transform: string | null = null;
      
      // If types don't match, apply appropriate transformation
      if (sourceField.type !== targetField.type) {
        if (sourceField.type === 'string' && targetField.type === 'number') {
          transform = 'parseFloat';
        } else if (sourceField.type === 'number' && targetField.type === 'string') {
          transform = 'toString';
        } else if (sourceField.type.startsWith('array') && targetField.type === 'string') {
          transform = 'splitFirst';
        }
      }
      
      mappings.push({
        id: `m${mappingCounter++}`,
        sourceField: sourceField.id,
        targetField: targetField.id,
        transform
      });
    }
  });
  
  return {
    sourceFields,
    targetFields,
    mappings
  };
};

/**
 * Extracts a value from a nested object using a path string
 * @param obj - The object to extract from
 * @param path - The path to the value (e.g. "user.address.city")
 * @returns The extracted value or null if not found
 */
export const extractValueFromPath = (obj: Record<string, any>, path: string): any => {
  if (!path) return obj;
  
  const parts = path.replace(/\[\]/g, "[0]").split('.');
  let current = obj;
  
  for (let part of parts) {
    const arrayMatch = part.match(/(.+)\[(\d+)\]$/);
    if (arrayMatch) {
      const [_, arrayName, index] = arrayMatch;
      if (!current[arrayName] || !current[arrayName][parseInt(index)]) {
        return null;
      }
      current = current[arrayName][parseInt(index)];
    } else {
      if (current[part] === undefined) {
        return null;
      }
      current = current[part];
    }
  }
  
  return current;
};
