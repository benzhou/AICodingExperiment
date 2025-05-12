# Transaction Matching Application Frontend

This is the frontend for the Transaction Matching Application, built with React, TypeScript, and Ant Design.

## Running the Application

### Important: Run from the Frontend Directory

You must run the application from the **frontend directory**, not the root project directory:

```bash
# Navigate to the frontend directory
cd /Users/benzhou/workspace/AICoding/frontend

# Install dependencies (if needed)
npm install

# Start the development server
npm start
```

Running from the root directory will fail with a "package.json not found" error.

## CSV Upload Troubleshooting

If you encounter issues when uploading CSV files:

1. **File Format**: Make sure your CSV file has headers in the first row and data starting from the second row
2. **File Type**: Save as "CSV (Comma delimited) (.csv)" if using Excel
3. **Delimiters**: The app supports comma, semicolon, and tab delimiters
4. **Analyze Tool**: Use the "Analyze CSV Format" button on the upload screen to diagnose issues
5. **Debug View**: Use the "Show Raw File Content" option to inspect your file

## Common Issues and Solutions

- **"No data has been parsed from your file"**: Check the file format and try using the Analyze CSV Format tool
- **Empty preview table**: Check if your CSV has proper headers and data rows
- **Column mapping not working**: Make sure your headers match expected patterns (date, amount, description, etc.)
- **Upload seems to hang**: Check browser console for errors (F12 > Console)

## Development Notes

For developers debugging CSV parsing issues:
- The file parsing code can be found in `src/components/DataSourceUpload.tsx`
- Set `NODE_ENV=development` to see debugging tools
- Check browser console for detailed parsing logs 