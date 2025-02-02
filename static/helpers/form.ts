type ParsedFormData = Record<string, Field>;
type Result = Record<string, string | string[]>
type Field = string | string[] | Result

function createNestedStructure(obj: Record<string, any>, keys: string[], value: Field): void {
  let currentObj = obj;
  for (let i = 0; i < keys.length - 1; i++) {
    const propertyName = keys[i];
    if (!currentObj[propertyName]) {
      // If the next key is a number, initialize an array, otherwise initialize an object
      currentObj[propertyName] = isNaN(Number(keys[i + 1])) ? {} : [];
    }
    currentObj = currentObj[propertyName];
  }

  const finalPropertyName = keys[keys.length - 1];
  if (Array.isArray(currentObj)) {
    // If the current object is an array, push the value(s)
    const index = parseInt(finalPropertyName, 10);
    if (currentObj[index] === undefined) {
      currentObj[index] = value;
    } else if (Array.isArray(currentObj[index])) {
      currentObj[index].push(...(Array.isArray(value) ? value : [value]));
    } else {
      currentObj[index] = [currentObj[index], ...(Array.isArray(value) ? value : [value])].flat();
    }
  } else {
    // If it's an object, assign the value to the property
    currentObj[finalPropertyName] = value;
  }
}

export async function parseForm(request: Request): Promise<Record<string, any>> {
  const formData = await request.formData();

  // Filter out files from formData
  const filteredFormData: ParsedFormData = {};
  for (const [key, value] of formData.entries()) {
    if (value instanceof File) {
      // Skip files
      continue;
    }
    // Convert to string or string array
    filteredFormData[key] = Array.isArray(value) ? Array.from(value) : value;
  }
  const transformedObj: Record<string, any> = {};

  for (const key in filteredFormData) {
    const keys = key.split(':');
    const finalValue = filteredFormData[key];
    createNestedStructure(transformedObj, keys, finalValue);
  }

  return transformedObj;
}
