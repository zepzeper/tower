import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../context/ThemeContext';
import FieldPanel from './mappings/FieldPanel';
import MappingPanel from './mappings/MappingPanel';
import MappingHeader from './mappings/MappingHeader';
import MappingPreview from './mappings/MappingPreview';
import MappingAlerts from './mappings/MappingAlerts';
import TransformPicker from './mappings/TransformPicker';
import {
  FieldDefinition,
  MappingDefinition,
  TransformationType,
  ConnectionDetails
} from '../../types';
import schemaService from '../../services/schemaService';
import * as schemaTransformer from '../../utils/schemaTransformer';
import { 
  ArrowLeftRight, 
  Plus, 
  Trash2, 
  Save, 
  RefreshCw, 
  ArrowRight, 
  Info, 
  AlertCircle,
  CheckCircle,
  ChevronDown,
  ChevronRight,
  Search,
  Settings,
  Eye,
  Download
} from 'lucide-react';

const FieldMapping: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme() || {};
  const { id: connectionId } = useParams<{ id: string }>();
  const isDarkMode = theme === 'dark';

  // States
  const [sourceFields, setSourceFields] = useState([]);
  const [targetFields, setTargetFields] = useState([]);
  const [mappings, setMappings] = useState([]);
  const [searchSource, setSearchSource] = useState('');
  const [searchTarget, setSearchTarget] = useState('');
  const [connectionDetails, setConnectionDetails] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [testData, setTestData] = useState(null);
  const [previewResult, setPreviewResult] = useState(null);
  const [activeTransformation, setActiveTransformation] = useState(null);
  const [showDragHelp, setShowDragHelp] = useState(true);
  const [dragSourceField, setDragSourceField] = useState(null);
  const [dragTargetField, setDragTargetField] = useState(null);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [mappingHistory, setMappingHistory] = useState([]);
  const [sourceFieldsExpanded, setSourceFieldsExpanded] = useState({});
  const [targetFieldsExpanded, setTargetFieldsExpanded] = useState({});
  const [loadingMapping, setIsLoadingMapping] = useState(false);

    useEffect(() => {
        const loadSchemas = async () => {
            setIsLoading(true);

            try {
                const connectionIdSafe = connectionId ?? 'demo-id'; // fallback just in case

                // Use service to get saved mappings and schemas
                const savedMapping = await schemaService.getSavedMapping(connectionIdSafe);
                const sourceType = savedMapping?.sourceType ?? 'woocommerce';
                const targetType = savedMapping?.targetType ?? 'brincr';

                const { sourceSchema, targetSchema } = await schemaService.getConnectionSchemas(sourceType, targetType);

                // Use schemaTransformer to flatten and prepare fields
                const { sourceFields, targetFields } = schemaTransformer.generateMappingData(sourceSchema, targetSchema);

                setSourceFields(sourceFields);
                setTargetFields(targetFields);

                // Use saved mappings if available
                if (savedMapping?.mappings?.length) {
                    setMappings(savedMapping.mappings);
                }

                setConnectionDetails({
                    id: connectionIdSafe,
                    name: `${sourceType} to ${targetType}`,
                    source: sourceType,
                    target: targetType,
                    status: 'active'
                });

            } catch (error) {
                console.error('Error loading connection schemas:', error);
            } finally {
                setIsLoading(false);
            }
        };

        loadSchemas();
    }, [connectionId]);

    // Available transformation functions
    const transformations = [
        { id: 'identity', name: 'No Transformation', description: 'Use value as is' },
        { id: 'parseFloat', name: 'Parse to Number', description: 'Convert string to number' },
        { id: 'toString', name: 'Convert to String', description: 'Convert value to string' },
        { id: 'trim', name: 'Trim', description: 'Remove whitespace from start and end' },
        { id: 'uppercase', name: 'Uppercase', description: 'Convert to UPPERCASE' },
        { id: 'lowercase', name: 'Lowercase', description: 'Convert to lowercase' },
        { id: 'round', name: 'Round', description: 'Round number to nearest integer' },
        { id: 'splitFirst', name: 'Split (First)', description: 'Get first item after splitting' },
        { id: 'concatenate', name: 'Concatenate', description: 'Join multiple values', isAdvanced: true },
        { id: 'conditional', name: 'Conditional', description: 'If-then logic for values', isAdvanced: true },
        { id: 'dateFormat', name: 'Format Date', description: 'Convert between date formats', isAdvanced: true },
    ];

    // Filter fields based on search
    const filteredSourceFields = sourceFields.filter(field => 
        field.name.toLowerCase().includes(searchSource.toLowerCase())
    );

    const filteredTargetFields = targetFields.filter(field => 
        field.name.toLowerCase().includes(searchTarget.toLowerCase())
    );

    // Drag and drop handlers
    const handleDragStart = (field, type) => {
        if (type === 'source') {
            setDragSourceField(field);
        } else {
            setDragTargetField(field);
        }
    };

    const handleDragEnd = () => {
        if (dragSourceField && dragTargetField) {
            addMapping(dragSourceField, dragTargetField);
        }

        setDragSourceField(null);
        setDragTargetField(null);
    };

    // Function to add a new mapping
    const addMapping = (sourceField, targetField) => {
        // First, check if the target field is already mapped
        const existingMapping = mappings.find(m => m.targetField === targetField.id);
        if (existingMapping) {
            // Save to history before changing
            setMappingHistory([...mappingHistory, [...mappings]]);

            // Update existing mapping
            setMappings(mappings.map(m => 
                m.id === existingMapping.id ? { ...m, sourceField: sourceField.id } : m
            ));
        } else {
            const newMapping = {
                id: `m${mappings.length + 1}`,
                sourceField: sourceField.id,
                targetField: targetField.id,
                transform: null
            };

            // Save current mappings to history before adding new one
            setMappingHistory([...mappingHistory, [...mappings]]);

            setMappings([...mappings, newMapping]);
        }

        // Hide the drag help after first mapping
        setShowDragHelp(false);
    };

    // Function to remove a mapping
    const removeMapping = (mappingId) => {
        // Save current mappings to history before removing
        setMappingHistory([...mappingHistory, [...mappings]]);

        setMappings(mappings.filter(m => m.id !== mappingId));
    };

    // Function to update a mapping's transformation
    const updateTransformation = (mappingId, transform) => {
        // Save current mappings to history before updating
        setMappingHistory([...mappingHistory, [...mappings]]);

        setMappings(mappings.map(m => 
            m.id === mappingId ? { ...m, transform } : m
        ));
        setActiveTransformation(null);
    };

    // Function to undo the last mapping change
    const undoMappingChange = () => {
        if (mappingHistory.length > 0) {
            const previousMappings = mappingHistory[mappingHistory.length - 1];
            setMappings(previousMappings);

            // Remove the last history entry
            setMappingHistory(mappingHistory.slice(0, -1));
        }
    };

    // Function to auto-map fields
    const autoMapFields = () => {
        setIsLoadingMapping(true)

        // Save current mappings to history before auto-mapping
        setMappingHistory([...mappingHistory, [...mappings]]);

        setTimeout(() => {
            // Simple auto-mapping algorithm
            const newMappings = [...mappings];

            // For each target field
            targetFields.forEach(targetField => {
                // Check if it's already mapped
                const alreadyMapped = mappings.some(m => m.targetField === targetField.id);
                if (alreadyMapped) return;

                // Try to find a matching source field
                const matchingSourceField = sourceFields.find(sf => 
                    // Exact name match
                    sf.name.toLowerCase() === targetField.name.toLowerCase() ||
                        // Or source contains target
                        sf.name.toLowerCase().includes(targetField.name.toLowerCase()) ||
                        // Or path ends with target name
                        sf.path.toLowerCase().endsWith(targetField.name.toLowerCase())
                );

                if (matchingSourceField) {
                    // Check if this source field is available
                    const sourceFieldAlreadyUsed = mappings.some(m => m.sourceField === matchingSourceField.id);

                    if (!sourceFieldAlreadyUsed) {
                        // Add the mapping
                        newMappings.push({
                            id: `m${newMappings.length + 1}`,
                            sourceField: matchingSourceField.id,
                            targetField: targetField.id,
                            transform: matchingSourceField.type !== targetField.type ? 'parseFloat' : null
                        });
                    }
                }
            });

            setMappings(newMappings);
            setIsLoadingMapping(false)

            // Run test mappings to show results
            testMappings();
        }, 800);
    };

    const testMappings = async () => {
        // Simulate API call for test transformation
        setIsLoading(true);

        try {
            const sourceData = await schemaService.getSchema(connectionDetails?.source ?? 'woocommerce');

            const mockResult: Record<string, any> = {};

            mappings.forEach(mapping => {
                const sourceField = getFieldById(mapping.sourceField, 'source');
                const targetField = getFieldById(mapping.targetField, 'target');

                if (!sourceField || !targetField) return;

                // Extract value from nested path
                let sourceValue;
                try {
                    sourceValue = extractValueFromPath(sourceData, sourceField.path);
                } catch {
                    sourceValue = null;
                }

                // Apply transformation
                if (mapping.transform && sourceValue !== null) {
                    switch (mapping.transform) {
                        case 'parseFloat':
                            sourceValue = parseFloat(sourceValue);
                            break;
                        case 'toString':
                            sourceValue = String(sourceValue);
                            break;
                        case 'trim':
                            sourceValue = sourceValue.trim();
                            break;
                        case 'uppercase':
                            sourceValue = sourceValue.toUpperCase();
                            break;
                        case 'lowercase':
                            sourceValue = sourceValue.toLowerCase();
                            break;
                        case 'round':
                            sourceValue = Math.round(parseFloat(sourceValue));
                            break;
                        case 'splitFirst':
                            if (Array.isArray(sourceValue)) {
                                sourceValue = sourceValue[0];
                            } else if (typeof sourceValue === 'string') {
                                sourceValue = sourceValue.split(' ')[0];
                            }
                            break;
                        // Add others as needed
                    }
                }

                mockResult[targetField.name] = sourceValue;
            });

            setTestData(sourceData);
            setPreviewResult(mockResult);
        } catch (error) {
            console.error('Error testing mappings:', error);
        } finally {
            setIsLoading(false);
        }
    };

    // Function to save mappings
    const saveMappings = async () => {
        setIsSaving(true);

        // Prepare the data to send to the server
        const mappingData = mappings.map(mapping => {
            const sourceField = getFieldById(mapping.sourceField, 'source');
            const targetField = getFieldById(mapping.targetField, 'target');

            return {
                id: mapping.id,
                sourcePath: sourceField?.path,
                targetPath: targetField?.path,
                transform: mapping.transform,
                sourceFieldId: mapping.sourceField,
                targetFieldId: mapping.targetField,
            };
        });

        const payload = {
            connectionId: connectionId!,
            sourceType: connectionDetails?.source ?? 'woocommerce',
            targetType: connectionDetails?.target ?? 'brincr',
            mappings,
            mappingData
        };

        await schemaService.saveMapping(connectionId!, payload);
        setSaveSuccess(true);
    };

    // Function to get field by ID
    const getFieldById = (fieldId, type) => {
        const fields = type === 'source' ? sourceFields : targetFields;
        return fields.find(f => f.id === fieldId);
    };

    // Function to get display transform name
    const getTransformName = (transformId) => {
        if (!transformId) return 'No Transformation';
        const transform = transformations.find(t => t.id === transformId);
        return transform ? transform.name : transformId;
    };

    // Function to toggle field expansion
    const toggleFieldExpansion = (fieldId, type) => {
        if (type === 'source') {
            setSourceFieldsExpanded({
                ...sourceFieldsExpanded,
                [fieldId]: !sourceFieldsExpanded[fieldId]
            });
        } else {
            setTargetFieldsExpanded({
                ...targetFieldsExpanded,
                [fieldId]: !targetFieldsExpanded[fieldId]
            });
        }
    };

    // Export mappings as JSON
    const exportMappings = () => {
        const mappingData = mappings.map(mapping => {
            const sourceField = getFieldById(mapping.sourceField, 'source');
            const targetField = getFieldById(mapping.targetField, 'target');

            return {
                sourcePath: sourceField?.path,
                targetPath: targetField?.path,
                transform: mapping.transform
            };
        });

        const dataStr = JSON.stringify(mappingData, null, 2);
        const dataUri = `data:application/json;charset=utf-8,${encodeURIComponent(dataStr)}`;

        const exportName = `${connectionDetails.source}-to-${connectionDetails.target}-mappings.json`;

        const linkElement = document.createElement('a');
        linkElement.setAttribute('href', dataUri);
        linkElement.setAttribute('download', exportName);
        linkElement.click();
    };

  return (
    <div className="p-6">
            <MappingHeader
                connectionName={connectionDetails?.name || ''}
                description={t('mappings.description')}
                isSaving={isSaving}
                isDarkMode={theme === 'dark'}
                t={t}
                onSave={saveMappings}
                onTest={testMappings}
                onExport={exportMappings}
            />


      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <FieldPanel
                    title={`${connectionDetails?.source} ${t('mappings.sourceFields')}`}
                    type="source"
                    fields={sourceFields}
                    mappings={mappings}
                    search={searchSource}
                    onSearchChange={setSearchSource}
                    onDragStart={handleDragStart}
                    onDragEnd={handleDragEnd}
                    expanded={sourceFieldsExpanded}
                    toggleExpansion={toggleFieldExpansion}
                    isDarkMode={theme === 'dark'}
                    getFieldById={getFieldById}
                    t={t}
                />

                <MappingPanel
                    mappings={mappings}
                    sourceFields={sourceFields}
                    targetFields={targetFields}
                    transformations={transformations}
                    activeTransformation={activeTransformation}
                    isDarkMode={theme === 'dark'}
                    showAdvanced={showAdvanced}
                    mappingHistory={mappingHistory}
                    t={t}
                    getFieldById={getFieldById}
                    getTransformName={getTransformName}
                    setActiveTransformation={setActiveTransformation}
                    updateTransformation={updateTransformation}
                    removeMapping={removeMapping}
                    undoMappingChange={undoMappingChange}
                    autoMapFields={autoMapFields}
                    setShowAdvanced={setShowAdvanced}
                    loadingAutoMapping={loadingMapping}
                />

                <FieldPanel
                    title={`${connectionDetails?.target} ${t('mappings.targetFields')}`}
                    type="target"
                    fields={targetFields}
                    mappings={mappings}
                    search={searchTarget}
                    onSearchChange={setSearchTarget}
                    onDragStart={handleDragStart}
                    onDragEnd={handleDragEnd}
                    expanded={targetFieldsExpanded}
                    toggleExpansion={toggleFieldExpansion}
                    isDarkMode={theme === 'dark'}
                    getFieldById={getFieldById}
                    t={t}
                    addMapping={addMapping}
                />
      </div>

            <MappingAlerts
                saveSuccess={saveSuccess}
                missingFields={targetFields
                    .filter(f => f.required && !mappings.some(m => m.targetField === f.id))
                    .map(f => f.name)}
                isDarkMode={theme === 'dark'}
                t={t}
            />
            <MappingPreview
                testData={testData}
                previewResult={previewResult}
                isDarkMode={theme === 'dark'}
                t={t}
            />
    </div>
  );
};

export default FieldMapping;
