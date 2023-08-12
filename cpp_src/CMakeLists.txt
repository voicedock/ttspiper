cmake_minimum_required(VERSION 3.14)
project(ttspiperlib
    VERSION 1.0.0
    DESCRIPTION "TTS piper lib"
    LANGUAGES C CXX)
    
#----------------------------------------------------------------------------------------------------------------------
# general settings and options
#----------------------------------------------------------------------------------------------------------------------

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

string(APPEND CMAKE_CXX_FLAGS " -Wall -Wextra -Wl,-rpath,'$ORIGIN'")
string(APPEND CMAKE_C_FLAGS " -Wall -Wextra")

include(GNUInstallDirs)

add_library(ttspiperlib SHARED
    src/ttspiperlib.cpp
    src/piper.cpp
)

set_target_properties(ttspiperlib PROPERTIES VERSION ${PROJECT_VERSION})
set_target_properties(ttspiperlib PROPERTIES SOVERSION 1)

#----------------------------------------------------------------------------------------------------------------------
# ttspiperlib dependencies
#----------------------------------------------------------------------------------------------------------------------
set(PIPER_PHONEMIZE_ROOTDIR ${CMAKE_CURRENT_LIST_DIR}/lib/${CMAKE_HOST_SYSTEM_NAME}-${CMAKE_HOST_SYSTEM_PROCESSOR}/piper_phonemize)
find_package(PkgConfig)
pkg_check_modules(SPDLOG REQUIRED spdlog)

target_link_libraries(ttspiperlib
                      piper_phonemize
                      espeak-ng
                      onnxruntime
                      pthread
                      ${SPDLOG_LIBRARIES}
                      -static-libgcc -static-libstdc++)

target_link_directories(ttspiperlib PRIVATE
                        ${PIPER_PHONEMIZE_ROOTDIR}/lib)

target_include_directories(ttspiperlib PRIVATE
                           ${PIPER_PHONEMIZE_ROOTDIR}/include
                           ${SPDLOG_INCLUDE_DIRS})

target_compile_options(ttspiperlib PRIVATE
                       ${SPDLOG_CFLAGS_OTHER})

#----------------------------------------------------------------------------------------------------------------------
# ttspiperlib sources
#----------------------------------------------------------------------------------------------------------------------
include_directories(include)
set_target_properties(ttspiperlib PROPERTIES PUBLIC_HEADER include/ttspiperlib/ttspiperlib.h)

#----------------------------------------------------------------------------------------------------------------------
# ttspiperlib target
#----------------------------------------------------------------------------------------------------------------------
install(TARGETS ttspiperlib
    LIBRARY DESTINATION ${CMAKE_INSTALL_LIBDIR}
    PUBLIC_HEADER DESTINATION ${CMAKE_INSTALL_INCLUDEDIR})