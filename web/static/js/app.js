class ImageProcessor {
    constructor() {
        this.uploadedFilename = null;
        this.initializeElements();
        this.setupEventListeners();
    }

    initializeElements() {
        // Upload elements
        this.uploadArea = document.getElementById('uploadArea');
        this.fileInput = document.getElementById('fileInput');
        this.uploadStatus = document.getElementById('uploadStatus');

        // Processing elements
        this.processingSection = document.getElementById('processingSection');
        this.resizeCheck = document.getElementById('resizeCheck');
        this.resizeControls = document.getElementById('resizeControls');
        this.grayscaleCheck = document.getElementById('grayscaleCheck');
        this.blurCheck = document.getElementById('blurCheck');
        this.blurControls = document.getElementById('blurControls');
        this.blurRadius = document.getElementById('blurRadius');
        this.blurValue = document.getElementById('blurValue');
        this.processBtn = document.getElementById('processBtn');
        this.resetBtn = document.getElementById('resetBtn');
        this.processStatus = document.getElementById('processStatus');

        // Result elements
        this.resultSection = document.getElementById('resultSection');
        this.downloadLink = document.getElementById('downloadLink');
        this.processAnotherBtn = document.getElementById('processAnotherBtn');

        // Loading element
        this.loading = document.getElementById('loading');
    }

    setupEventListeners() {
        // Upload area events
        this.uploadArea.addEventListener('click', () => this.fileInput.click());
        this.uploadArea.addEventListener('dragover', (e) => this.handleDragOver(e));
        this.uploadArea.addEventListener('dragleave', (e) => this.handleDragLeave(e));
        this.uploadArea.addEventListener('drop', (e) => this.handleDrop(e));

        // File input change
        this.fileInput.addEventListener('change', (e) => this.handleFileSelect(e));

        // Processing controls
        this.resizeCheck.addEventListener('change', () => this.toggleResizeControls());
        this.blurCheck.addEventListener('change', () => this.toggleBlurControls());
        this.blurRadius.addEventListener('input', () => this.updateBlurValue());

        // Action buttons
        this.processBtn.addEventListener('click', () => this.processImage());
        this.resetBtn.addEventListener('click', () => this.reset());
        this.processAnotherBtn.addEventListener('click', () => this.reset());
    }

    handleDragOver(e) {
        e.preventDefault();
        this.uploadArea.classList.add('dragover');
    }

    handleDragLeave(e) {
        e.preventDefault();
        this.uploadArea.classList.remove('dragover');
    }

    handleDrop(e) {
        e.preventDefault();
        this.uploadArea.classList.remove('dragover');
        
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            this.uploadFile(files[0]);
        }
    }

    handleFileSelect(e) {
        const file = e.target.files[0];
        if (file) {
            this.uploadFile(file);
        }
    }

    async uploadFile(file) {
        // Validate file
        if (!this.validateFile(file)) {
            return;
        }

        this.showLoading();

        const formData = new FormData();
        formData.append('image', file);

        try {
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();

            if (response.ok) {
                this.uploadedFilename = result.filename;
                this.showUploadSuccess(result.message);
                this.showProcessingSection();
            } else {
                this.showUploadError(result.message || 'Upload failed');
            }
        } catch (error) {
            this.showUploadError('Network error occurred');
        } finally {
            this.hideLoading();
        }
    }

    validateFile(file) {
        const maxSize = 10 * 1024 * 1024; // 10MB
        const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif'];

        if (file.size > maxSize) {
            this.showUploadError('File size must be less than 10MB');
            return false;
        }

        if (!allowedTypes.includes(file.type)) {
            this.showUploadError('Only JPG, PNG, and GIF files are allowed');
            return false;
        }

        return true;
    }

    async processImage() {
        if (!this.uploadedFilename) {
            this.showProcessError('No file uploaded');
            return;
        }

        const operations = this.buildOperations();
        if (operations.length === 0) {
            this.showProcessError('Please select at least one processing option');
            return;
        }

        this.showLoading();

        const requestData = {
            filename: this.uploadedFilename,
            operations: operations
        };

        try {
            const response = await fetch('/process', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            });

            const result = await response.json();

            if (response.ok) {
                this.showProcessSuccess(result.message);
                this.showResultSection(result.download_url);
            } else {
                this.showProcessError(result.message || 'Processing failed');
            }
        } catch (error) {
            this.showProcessError('Network error occurred');
        } finally {
            this.hideLoading();
        }
    }

    buildOperations() {
        const operations = [];

        // Add resize operation
        if (this.resizeCheck.checked) {
            const width = parseInt(document.getElementById('width').value);
            const height = parseInt(document.getElementById('height').value);
            
            if (width > 0 && height > 0) {
                operations.push({
                    type: 'resize',
                    params: { width, height }
                });
            }
        }

        // Add grayscale operation
        if (this.grayscaleCheck.checked) {
            operations.push({
                type: 'grayscale'
            });
        }

        // Add blur operation
        if (this.blurCheck.checked) {
            const radius = parseFloat(this.blurRadius.value);
            operations.push({
                type: 'blur',
                params: { radius }
            });
        }

        return operations;
    }

    toggleResizeControls() {
        if (this.resizeCheck.checked) {
            this.resizeControls.classList.remove('hidden');
        } else {
            this.resizeControls.classList.add('hidden');
        }
    }

    toggleBlurControls() {
        if (this.blurCheck.checked) {
            this.blurControls.classList.remove('hidden');
        } else {
            this.blurControls.classList.add('hidden');
        }
    }

    updateBlurValue() {
        this.blurValue.textContent = this.blurRadius.value;
    }

    showUploadSuccess(message) {
        this.uploadStatus.textContent = message;
        this.uploadStatus.className = 'status success';
        this.uploadStatus.classList.remove('hidden');
    }

    showUploadError(message) {
        this.uploadStatus.textContent = message;
        this.uploadStatus.className = 'status error';
        this.uploadStatus.classList.remove('hidden');
    }

    showProcessSuccess(message) {
        this.processStatus.textContent = message;
        this.processStatus.className = 'status success';
        this.processStatus.classList.remove('hidden');
    }

    showProcessError(message) {
        this.processStatus.textContent = message;
        this.processStatus.className = 'status error';
        this.processStatus.classList.remove('hidden');
    }

    showProcessingSection() {
        this.processingSection.classList.remove('hidden');
        this.resultSection.classList.add('hidden');
    }

    showResultSection(downloadUrl) {
        this.resultSection.classList.remove('hidden');
        this.processingSection.classList.add('hidden');
        this.downloadLink.href = downloadUrl;
    }

    showLoading() {
        this.loading.classList.remove('hidden');
    }

    hideLoading() {
        this.loading.classList.add('hidden');
    }

    reset() {
        // Reset form state
        this.uploadedFilename = null;
        this.fileInput.value = '';
        
        // Reset checkboxes
        this.resizeCheck.checked = false;
        this.grayscaleCheck.checked = false;
        this.blurCheck.checked = false;
        
        // Hide controls
        this.resizeControls.classList.add('hidden');
        this.blurControls.classList.add('hidden');
        
        // Reset values
        document.getElementById('width').value = '800';
        document.getElementById('height').value = '600';
        this.blurRadius.value = '2';
        this.updateBlurValue();
        
        // Hide sections and status messages
        this.processingSection.classList.add('hidden');
        this.resultSection.classList.add('hidden');
        this.uploadStatus.classList.add('hidden');
        this.processStatus.classList.add('hidden');
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new ImageProcessor();
});