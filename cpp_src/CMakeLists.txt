cmake_minimum_required(VERSION 3.14)
project(ttspiperlib
    VERSION 1.0.0
    DESCRIPTION "TTS piper lib"
    LANGUAGES CXX)
    
#----------------------------------------------------------------------------------------------------------------------
# general settings and options
#----------------------------------------------------------------------------------------------------------------------

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(GNUInstallDirs)

add_library(ttspiperlib SHARED
    src/ttspiperlib.cpp
)

set_target_properties(ttspiperlib PROPERTIES VERSION ${PROJECT_VERSION})
set_target_properties(ttspiperlib PROPERTIES SOVERSION 1)

#----------------------------------------------------------------------------------------------------------------------
# ttspiperlib dependencies
#----------------------------------------------------------------------------------------------------------------------
find_package(PkgConfig)
pkg_check_modules(ESPEAK_NG REQUIRED espeak-ng<2)
set(ONNXRUNTIME_ROOTDIR ${CMAKE_CURRENT_LIST_DIR}/lib/${CMAKE_HOST_SYSTEM_NAME}-${CMAKE_HOST_SYSTEM_PROCESSOR})
target_link_libraries(ttspiperlib
                      onnxruntime
                      pthread
                      -static-libgcc -static-libstdc++
                      ${ESPEAK_NG_LIBRARIES})

target_link_directories(ttspiperlib PRIVATE
                        ${ESPEAK_NG_LIBRARY_DIRS}
                        ${ONNXRUNTIME_ROOTDIR}/lib)

target_include_directories(ttspiperlib PRIVATE
                           ${ONNXRUNTIME_ROOTDIR}/include
                           ${ESPEAK_NG_INCLUDE_DIRS})

target_compile_options(ttspiperlib PUBLIC
                       ${ESPEAK_NG_CFLAGS_OTHER})

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