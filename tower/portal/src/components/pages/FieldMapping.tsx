import React, { useEffect, useState, useRef, useCallback } from 'react';
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
  const lastTestedMappingKey = useRef('');

  useEffect(() => {
    // Only run when mapping changes and not during loading
    if (!loadingMapping && mappings.length > 0) {
      // Use a debounce timer to avoid multiple API calls
      const debounceTimer = setTimeout(() => {
        testMappings();
      }, 1000); // Wait 1 second after changes before testing

      // Clean up the timer if mappings change again
      return () => clearTimeout(debounceTimer);
    }
  }, [mappings, loadingMapping]);

  useEffect(() => {
    const loadSchemas = async () => {
      // Don't set loading if we already have data
      if (sourceFields.length === 0 || targetFields.length === 0) {
        setIsLoading(true);
      }

      try {
        const connectionIdSafe = connectionId ?? 'demo-id';

        // Check if we already have connection details to avoid unnecessary calls
        if (!connectionDetails || connectionDetails.id !== connectionIdSafe) {
          // Use service to get saved mappings and schemas
          const sourceType = 'woocommerce';
          const targetType = 'brincr';

          // Only fetch schemas if we don't have them already
          if (sourceFields.length === 0 || targetFields.length === 0) {
            const { sourceFields: source, targetFields: target, mappings: mappings } =
              await schemaService.getConnectionSchemas(sourceType, targetType);

            setSourceFields(source);
            setTargetFields(target);

            // Use saved mappings if available and we don't already have mappings
            if (mappings?.length) {
              setMappings(mappings);
            }
          }

          setConnectionDetails({
            id: connectionIdSafe,
            name: `${sourceType} to ${targetType}`,
            source: sourceType,
            target: targetType,
            status: 'active'
          });
        }
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
    setIsLoadingMapping(true);
    setMappingHistory([...mappingHistory, [...mappings]]);

    setTimeout(() => {
      // Collect all new mappings at once to avoid multiple state updates
      const newMappings = [...mappings];
      let mappingsAdded = false;

      targetFields.forEach(targetField => {
        const alreadyMapped = mappings.some(m => m.targetField === targetField.id);
        if (alreadyMapped) return;

        const matchingSourceField = sourceFields.find(sf =>
          sf.name.toLowerCase() === targetField.name.toLowerCase() ||
          sf.name.toLowerCase().includes(targetField.name.toLowerCase()) ||
          sf.path.toLowerCase().endsWith(targetField.name.toLowerCase())
        );

        if (matchingSourceField) {
          const sourceFieldAlreadyUsed = newMappings.some(m => m.sourceField === matchingSourceField.id);

          if (!sourceFieldAlreadyUsed) {
            newMappings.push({
              id: `m${newMappings.length + 1}`,
              sourceField: matchingSourceField.id,
              targetField: targetField.id,
              transform: matchingSourceField.type !== targetField.type ? 'parseFloat' : null
            });
            mappingsAdded = true;
          }
        }
      });

      // Only update state if we actually added mappings
      if (mappingsAdded) {
        setMappings(newMappings);
      }

      setIsLoadingMapping(false);
    }, 800);
  };

  // Function to get field by ID - Move this up before testMappings
  const getFieldById = (fieldId, type) => {
    const fields = type === 'source' ? sourceFields : targetFields;
    return fields.find(f => f.id === fieldId);
  };

  // Optimize testMappings to cache results and avoid redundant API calls
  const testMappings = useCallback(async () => {
    // Avoid testing if there are no mappings
    if (mappings.length === 0) {
      setTestData(null);
      setPreviewResult(null);
      return;
    }

    setIsLoading(true);

    try {
      // Generate a hash/key for the current mappings to check if we need to retest
      const mappingKey = JSON.stringify(mappings);

      // If we already tested these exact mappings, use cached results
      if (mappingKey === lastTestedMappingKey.current) {
        setIsLoading(false);
        return;
      }

      // Prepare mapping metadata for API call
      const mappingMetadata = mappings.map(mapping => {
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

      // Use the schema service
      const result = await schemaService.testMapping(
        connectionId ?? 'demo-id',
        connectionDetails?.source ?? 'woocommerce',
        connectionDetails?.target ?? 'brincr',
        mappingMetadata
      );

      // Set test data (source) and preview result (transformed)
      setTestData(result.sourceData);
      setPreviewResult(result.transformedData);

      // Store this mapping key as last tested
      lastTestedMappingKey.current = mappingKey;

    } catch (error) {
      console.error('Error testing mappings:', error);
    } finally {
      setIsLoading(false);
    }
  }, [mappings, connectionId, connectionDetails, sourceFields, targetFields]);

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
        operation: "products",
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
